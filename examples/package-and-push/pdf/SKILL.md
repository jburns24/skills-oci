---
name: pdf
version: 1.0.0
description: Process PDF files — merge, split, extract text and tables, create PDFs, and handle scanned documents with OCR.
license: Apache-2.0
compatibility: |
  Python libraries: pypdf, pdfplumber, reportlab, pytesseract, pdf2image.
  CLI tools: pdftotext (poppler-utils), qpdf, pdftk, pdfimages.
  Install Python deps: pip install pypdf pdfplumber reportlab pytesseract pdf2image pandas
metadata:
  category: document-processing
  tags: [pdf, documents, text-extraction, ocr, merge, split]
---

# PDF Processing

This skill provides capabilities for processing PDF files using Python libraries and command-line tools.

## Quick Start

```python
from pypdf import PdfReader, PdfWriter

# Read a PDF
reader = PdfReader("document.pdf")
print(f"Pages: {len(reader.pages)}")

# Extract text
text = ""
for page in reader.pages:
    text += page.extract_text()
```

## Core Operations

### Merge PDFs

```python
from pypdf import PdfWriter, PdfReader

writer = PdfWriter()
for pdf_file in ["doc1.pdf", "doc2.pdf", "doc3.pdf"]:
    reader = PdfReader(pdf_file)
    for page in reader.pages:
        writer.add_page(page)

with open("merged.pdf", "wb") as output:
    writer.write(output)
```

### Split PDF

```python
reader = PdfReader("input.pdf")
for i, page in enumerate(reader.pages):
    writer = PdfWriter()
    writer.add_page(page)
    with open(f"page_{i+1}.pdf", "wb") as output:
        writer.write(output)
```

### Extract Text with Layout

```python
import pdfplumber

with pdfplumber.open("document.pdf") as pdf:
    for page in pdf.pages:
        text = page.extract_text()
        print(text)
```

### Extract Tables

```python
import pdfplumber
import pandas as pd

with pdfplumber.open("document.pdf") as pdf:
    all_tables = []
    for page in pdf.pages:
        tables = page.extract_tables()
        for table in tables:
            if table:
                df = pd.DataFrame(table[1:], columns=table[0])
                all_tables.append(df)

if all_tables:
    combined_df = pd.concat(all_tables, ignore_index=True)
    combined_df.to_excel("extracted_tables.xlsx", index=False)
```

### Create a PDF

```python
from reportlab.lib.pagesizes import letter
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer
from reportlab.lib.styles import getSampleStyleSheet

doc = SimpleDocTemplate("report.pdf", pagesize=letter)
styles = getSampleStyleSheet()
story = [
    Paragraph("Report Title", styles['Title']),
    Spacer(1, 12),
    Paragraph("Body content goes here.", styles['Normal']),
]
doc.build(story)
```

### OCR Scanned PDFs

```python
import pytesseract
from pdf2image import convert_from_path

images = convert_from_path('scanned.pdf')
text = ""
for i, image in enumerate(images):
    text += f"Page {i+1}:\n"
    text += pytesseract.image_to_string(image)
    text += "\n\n"
print(text)
```

## Command-Line Tools

```bash
# Extract text (preserving layout)
pdftotext -layout input.pdf output.txt

# Merge PDFs
qpdf --empty --pages file1.pdf file2.pdf -- merged.pdf

# Split pages 1-5
qpdf input.pdf --pages . 1-5 -- pages1-5.pdf

# Rotate page 1 by 90 degrees
qpdf input.pdf output.pdf --rotate=+90:1

# Remove password
qpdf --password=mypassword --decrypt encrypted.pdf decrypted.pdf

# Extract images
pdfimages -j input.pdf output_prefix
```

## Password Protection

```python
from pypdf import PdfReader, PdfWriter

reader = PdfReader("input.pdf")
writer = PdfWriter()
for page in reader.pages:
    writer.add_page(page)

writer.encrypt("userpassword", "ownerpassword")
with open("encrypted.pdf", "wb") as output:
    writer.write(output)
```

## Quick Reference

| Task | Best Tool | Notes |
|------|-----------|-------|
| Merge PDFs | pypdf / qpdf | |
| Split PDFs | pypdf / qpdf | |
| Extract text | pdfplumber | Layout-aware |
| Extract tables | pdfplumber + pandas | Export to xlsx |
| Create PDFs | reportlab | Canvas or Platypus |
| OCR scanned PDFs | pytesseract | Convert to image first |
| Password protect | pypdf | User + owner passwords |
| Extract images | pdfimages | poppler-utils |

> **Note:** This skill is based on the `pdf` skill from [anthropics/skills](https://github.com/anthropics/skills), one of the most popular skills on [skills.sh](https://skills.sh). It has been repackaged as an OCI artifact for distribution via `skills-oci`.
