# PDF Size Limits and Performance

## Upload Limits

| Component | Default Limit | Notes |
|-----------|---------------|-------|
| FastAPI | ~100MB | Default upload size |
| PyMuPDF | No hard limit | Handles 1000+ pages |
| System RAM | Varies | PDF loads into memory |

## Performance by PDF Type

### Native Text PDFs (Digital)

| Pages | File Size | Processing Time |
|-------|-----------|-----------------|
| 1-10 | < 1MB | < 100ms |
| 50-100 | 5-10MB | < 1 second |
| 500+ | 50MB+ | 2-5 seconds |

Native text extraction is fast - limited only by file I/O.

### Scanned/OCR PDFs

| Pages | Processing Time | Notes |
|-------|-----------------|-------|
| 1 | 2-3 seconds | 300 DPI rendering |
| 10 | 20-30 seconds | CPU intensive |
| 50 | 2-3 minutes | Consider batching |
| 100+ | 5+ minutes | May need timeout adjustment |

OCR is CPU-bound. Each page is rendered at 300 DPI then processed by Tesseract.

## Memory Usage

- **Native PDF**: ~2-3x file size in RAM
- **OCR PDF**: ~10x file size in RAM (image rendering)

Example: A 10MB scanned PDF may use ~100MB RAM during OCR processing.

## Increasing Upload Limit

To allow larger files, modify `simple_pdf_extract_api.py`:

```python
from fastapi import FastAPI

app = FastAPI()

# Add this middleware for larger uploads
from starlette.middleware.base import BaseHTTPMiddleware

class MaxBodySizeMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, max_body_size: int):
        super().__init__(app)
        self.max_body_size = max_body_size

    async def dispatch(self, request, call_next):
        return await call_next(request)

# 500MB limit
app.add_middleware(MaxBodySizeMiddleware, max_body_size=500_000_000)
```

Or set via uvicorn:

```bash
uvicorn simple_pdf_extract_api:app --host 0.0.0.0 --port 8000 --limit-max-request-size 500000000
```

## Timeout Considerations

For large OCR jobs, increase the timeout:

```bash
uvicorn simple_pdf_extract_api:app --host 0.0.0.0 --port 8000 --timeout-keep-alive 300
```

## Recommendations

| Use Case | Recommendation |
|----------|----------------|
| Small PDFs (< 10 pages) | Default settings work fine |
| Large native PDFs | Increase upload limit if needed |
| Large scanned PDFs | Process in batches, increase timeout |
| Production deployment | Add request queue for OCR jobs |

## Testing Large Files

Before production, test with:
1. A 50+ page native text PDF
2. A 20+ page scanned PDF
3. A mixed content PDF with images

Monitor memory and processing time to establish baselines.
