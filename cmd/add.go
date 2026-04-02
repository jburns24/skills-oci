package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/salaboy/skills-oci/pkg/oci"
	"github.com/salaboy/skills-oci/pkg/skill"
	"github.com/salaboy/skills-oci/pkg/tui/add"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Install a skill from an OCI registry",
		Long:  "Pulls a skill artifact from a remote container registry, extracts it to the skills directory, and updates skills.json and skills.lock.json.",
		Example: `  # Install a skill from GHCR
  skills-oci add --ref ghcr.io/myorg/skills/my-skill:1.0.0

  # Install from a local registry
  skills-oci add --ref localhost:5000/my-skill:1.0.0 --plain-http

  # Install to .claude/skills instead of .agents/skills
  skills-oci add --ref ghcr.io/myorg/skills/my-skill:1.0.0 --claude

  # Install to a custom directory
  skills-oci add --ref ghcr.io/myorg/skills/my-skill:1.0.0 --output ./custom/path`,
		RunE: runAdd,
	}

	cmd.Flags().String("ref", "", "Full OCI reference (e.g., ghcr.io/org/skills/my-skill:1.0.0)")
	cmd.Flags().String("output", "", "Output directory for skill extraction (overrides default)")
	cmd.Flags().String("project-dir", ".", "Project directory containing skills.json and skills.lock.json")

	_ = cmd.MarkFlagRequired("ref")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	ref, _ := cmd.Flags().GetString("ref")
	output, _ := cmd.Flags().GetString("output")
	projectDir, _ := cmd.Flags().GetString("project-dir")
	plain, _ := cmd.Flags().GetBool("plain")
	plainHTTP, _ := cmd.Flags().GetBool("plain-http")
	skillsDir := resolveSkillsDir(cmd)

	// If no explicit output, use the resolved skills dir relative to project dir
	if output == "" {
		output = filepath.Join(projectDir, skillsDir)
	}

	if plain {
		return runAddPlain(ref, output, projectDir, skillsDir, plainHTTP)
	}

	m := add.NewModel(ref, output, projectDir, skillsDir, plainHTTP)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	if fm, ok := finalModel.(add.Model); ok {
		if fm.Err() != nil {
			return fm.Err()
		}
	}

	return nil
}

func runAddPlain(ref, output, projectDir, skillsDir string, plainHTTP bool) error {
	result, err := oci.Pull(context.Background(), oci.PullOptions{
		Reference: ref,
		OutputDir: output,
		PlainHTTP: plainHTTP,
		OnStatus: func(phase string) {
			fmt.Printf("  %s\n", phase)
		},
	})
	if err != nil {
		return err
	}

	// Update skills.json
	fmt.Println("  Updating skills.json")
	if err := updateManifest(projectDir, skillsDir, result); err != nil {
		return fmt.Errorf("updating skills.json: %w", err)
	}

	// Update skills.lock.json
	fmt.Println("  Updating skills.lock.json")
	if err := updateLockFile(projectDir, skillsDir, result); err != nil {
		return fmt.Errorf("updating skills.lock.json: %w", err)
	}

	fmt.Printf("\nSuccessfully installed!\n")
	fmt.Printf("  Name:      %s\n", result.Name)
	fmt.Printf("  Version:   %s\n", result.Version)
	fmt.Printf("  Digest:    %s\n", result.Digest)
	fmt.Printf("  Extracted: %s\n", result.ExtractTo)
	return nil
}

// updateManifest loads skills.json, adds/updates the skill entry, and saves it.
func updateManifest(projectDir, skillsDir string, result *oci.PullResult) error {
	m, err := skill.LoadManifest(projectDir)
	if err != nil {
		return err
	}

	skill.AddToManifest(m, result.Name, result.Source(), result.Version)

	return skill.SaveManifest(projectDir, m)
}

// updateLockFile loads skills.lock.json, adds/updates the skill entry, and saves it.
func updateLockFile(projectDir, skillsDir string, result *oci.PullResult) error {
	l, err := skill.LoadLock(projectDir)
	if err != nil {
		return err
	}

	extractPath := filepath.Join(skillsDir, result.Name)

	entry := skill.LockedSkill{
		Name: result.Name,
		Path: extractPath,
		Source: skill.LockSource{
			Registry:   result.Registry,
			Repository: result.Repository,
			Tag:        result.Tag,
			Digest:     result.Digest,
			Ref:        result.FullRef(),
		},
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
	}

	skill.AddToLock(l, entry)

	return skill.SaveLock(projectDir, l)
}
