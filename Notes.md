 Known Issues and Limitations
OCR Performance
Scanned PDFs are slower since OCR has to chew through each page - couple seconds per page, so big scanned docs take a bit. Native text PDFs are basically instant.
OCR Accuracy
Accuracy depends on scan quality. Clean scans are pretty good (90%+), but grainy or skewed stuff might have some weird characters or missed words. Expect 70-90% on rough scans.
Text in Images/Logos
Text baked into images or logos won't be picked up unless you use force_ocr=true. That tells it to OCR the whole page instead of just relying on native text extraction.
Multi-Column Layouts
Multi-column layouts can get jumbled since it reads left to right across the whole page. Tables and side-by-side content might not come out in the order you expect.
Memory Usage
OCR is hungry - figure 10x the file size in RAM while it's processing. A 10MB scanned PDF could use around 100MB during extraction.
Upload Limit
Default upload cap is around 100MB. Can be increased in the uvicorn config if needed:
uvicorn simple_pdf_extract_api:app --host 0.0.0.0 --port 8000 --limit-max-request-size 500000000

Timeouts on Large Files
Big scanned PDFs might hit timeout limits. For 50+ page scanned docs, consider increasing the timeout:
uvicorn simple_pdf_extract_api:app --host 0.0.0.0 --port 8000 --timeout-keep-alive 300


