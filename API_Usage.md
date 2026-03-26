# API Usage Guide

## Endpoint
http://<server-ip>:8000
## Health Check
```bash
curl http://<server-ip>:8000/health
```

Response:
```json
{"status":"healthy","timestamp":"2026-03-26T09:12:40.037578","ocr_available":true}
```

## Extract PDF (Standard)
```bash
curl -X POST -F "file=@your-document.pdf" http://<server-ip>:8000/extract
```

## Extract PDF with Force OCR

Use this to capture text from logos, images, and embedded graphics:
```bash
curl -X POST -F "file=@your-document.pdf" "http://<server-ip>:8000/extract?force_ocr=true"
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| file | file | required | PDF file to upload |
| use_ocr | bool | true | Enable OCR for scanned pages |
| force_ocr | bool | false | OCR entire page to capture text in images |

## Example Response
```json
{
  "success": true,
  "extraction_info": {
    "source_file": "document.pdf",
    "ocr_available": true,
    "force_ocr": false
  },
  "metadata": {
    "title": "Document Title",
    "author": "Author Name",
    "page_count": 5
  },
  "pages": [
    {
      "page_number": 1,
      "text": "Extracted text content...",
      "text_length": 1500,
      "extraction_method": "native",
      "has_text": true
    }
  ],
  "summary": {
    "total_pages": 5,
    "total_characters": 7500,
    "pages_with_text": 5,
    "pages_ocr_required": 0
  }
}
```

## Extraction Methods

| Method | Meaning |
|--------|---------|
| native | Text extracted directly from PDF |
| ocr | Text extracted via Tesseract OCR |
| native+ocr | Both native and OCR text (when force_ocr=true) |

## Save Output to File
```bash
curl -X POST -F "file=@document.pdf" http://<server-ip>:8000/extract -o output.json
```

## Python Example
```python
import requests

url = "http://<server-ip>:8000/extract"
files = {"file": open("document.pdf", "rb")}
params = {"force_ocr": "false"}

response = requests.post(url, files=files, params=params)
data = response.json()

print(f"Pages: {data['summary']['total_pages']}")
print(f"Characters: {data['summary']['total_characters']}")
```
