# PDF Extraction API

A FastAPI-based service for extracting text, metadata, and structure from PDF files using PyMuPDF with Tesseract OCR fallback.

## Features

- **Native text extraction** - Fast, accurate extraction from digital PDFs
- **OCR fallback** - Automatic Tesseract OCR for scanned/image-based PDFs
- **Force OCR mode** - Extract text from logos, images, and embedded graphics
- **Structured JSON output** - Metadata, page text, images, and links
- **REST API** - Simple HTTP endpoints for integration

## Quick Start

### Prerequisites

- Python 3.10+
- Tesseract OCR 5.x

### Installation

sudo apt update && sudo apt install -y tesseract-ocr tesseract-ocr-eng

python3 -m venv .venv
source .venv/bin/activate

pip install -r requirements.txt

### Run the API

uvicorn simple_pdf_extract_api:app --host 0.0.0.0 --port 8000

### Test the API

curl http://127.0.0.1:8000/health

curl -X POST -F "file=@document.pdf" http://127.0.0.1:8000/extract

curl -X POST -F "file=@document.pdf" "http://127.0.0.1:8000/extract?force_ocr=true"

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| / | GET | API info and status |
| /health | GET | Health check |
| /extract | POST | Extract PDF content |

### /extract Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| file | file | required | PDF file to process |
| use_ocr | bool | tru
