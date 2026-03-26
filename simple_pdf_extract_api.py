#!/usr/bin/env python3
"""
simple_pdf_extract_api.py - with force_ocr support
"""

import io
from datetime import datetime
from typing import Optional

from fastapi import FastAPI, File, UploadFile, HTTPException, Query
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware

try:
    import fitz
except ImportError:
    raise ImportError("PyMuPDF not installed. Run: pip install PyMuPDF")

try:
    import pytesseract
    from PIL import Image
    OCR_AVAILABLE = True
except ImportError:
    OCR_AVAILABLE = False

app = FastAPI(title="PDF Extraction API", version="1.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

def extract_pdf_metadata(doc: fitz.Document) -> dict:
    meta = doc.metadata or {}
    return {
        "title": meta.get("title", "") or "",
        "author": meta.get("author", "") or "",
        "subject": meta.get("subject", "") or "",
        "keywords": meta.get("keywords", "") or "",
        "creator": meta.get("creator", "") or "",
        "producer": meta.get("producer", "") or "",
        "creation_date": meta.get("creationDate", "") or "",
        "modification_date": meta.get("modDate", "") or "",
        "page_count": doc.page_count,
        "is_encrypted": doc.is_encrypted,
    }

def extract_page_text(page: fitz.Page, use_ocr: bool = False, force_ocr: bool = False) -> dict:
    native_text = page.get_text("text").strip()
    ocr_text = ""
    extraction_method = "native"
    
    if force_ocr and OCR_AVAILABLE:
        try:
            mat = fitz.Matrix(300/72, 300/72)
            pix = page.get_pixmap(matrix=mat)
            img_data = pix.tobytes("png")
            img = Image.open(io.BytesIO(img_data))
            ocr_text = pytesseract.image_to_string(img).strip()
        except Exception:
            pass
        
        if native_text and ocr_text:
            text = f"=== NATIVE TEXT ===\n{native_text}\n\n=== OCR TEXT (from images) ===\n{ocr_text}"
            extraction_method = "native+ocr"
        elif ocr_text:
            text = ocr_text
            extraction_method = "force_ocr"
        else:
            text = native_text
    elif not native_text and use_ocr and OCR_AVAILABLE:
        try:
            mat = fitz.Matrix(300/72, 300/72)
            pix = page.get_pixmap(matrix=mat)
            img_data = pix.tobytes("png")
            img = Image.open(io.BytesIO(img_data))
            text = pytesseract.image_to_string(img).strip()
            extraction_method = "ocr"
        except Exception:
            text = ""
            extraction_method = "ocr_failed"
    else:
        text = native_text
    
    return {
        "page_number": page.number + 1,
        "text": text,
        "text_length": len(text),
        "extraction_method": extraction_method,
        "has_text": bool(text),
    }

def extract_page_images(page: fitz.Page) -> list:
    images = []
    try:
        for idx, img in enumerate(page.get_images(full=True)):
            images.append({"image_index": idx + 1, "xref": img[0], "width": img[2], "height": img[3]})
    except Exception:
        pass
    return images

def extract_page_links(page: fitz.Page) -> list:
    links = []
    try:
        for link in page.get_links():
            link_info = {"kind": link.get("kind", 0), "uri": link.get("uri", ""), "page": link.get("page", -1)}
            if link_info["uri"] or link_info["page"] >= 0:
                links.append(link_info)
    except Exception:
        pass
    return links

def process_pdf(pdf_bytes: bytes, filename: str, use_ocr: bool = True, force_ocr: bool = False) -> dict:
    try:
        doc = fitz.open(stream=pdf_bytes, filetype="pdf")
    except Exception as e:
        raise ValueError(f"Failed to open PDF: {str(e)}")
    
    result = {
        "success": True,
        "extraction_info": {
            "source_file": filename,
            "extraction_timestamp": datetime.now().isoformat(),
            "pymupdf_version": fitz.version[0],
            "ocr_available": OCR_AVAILABLE,
            "ocr_used": use_ocr and OCR_AVAILABLE,
            "force_ocr": force_ocr,
        },
        "metadata": extract_pdf_metadata(doc),
        "pages": [],
        "summary": {
            "total_pages": doc.page_count,
            "total_characters": 0,
            "total_images": 0,
            "total_links": 0,
            "pages_with_text": 0,
            "pages_ocr_required": 0,
        }
    }
    
    for page_num in range(doc.page_count):
        page = doc[page_num]
        page_data = extract_page_text(page, use_ocr=use_ocr, force_ocr=force_ocr)
        page_data["images"] = extract_page_images(page)
        page_data["links"] = extract_page_links(page)
        page_data["width"] = page.rect.width
        page_data["height"] = page.rect.height
        result["pages"].append(page_data)
        result["summary"]["total_characters"] += page_data["text_length"]
        result["summary"]["total_images"] += len(page_data["images"])
        result["summary"]["total_links"] += len(page_data["links"])
        if page_data["has_text"]:
            result["summary"]["pages_with_text"] += 1
        if page_data["extraction_method"] in ["ocr", "force_ocr", "native+ocr"]:
            result["summary"]["pages_ocr_required"] += 1
    
    doc.close()
    return result

@app.get("/")
async def root():
    return {"service": "PDF Extraction API", "version": "1.1.0", "status": "running", "ocr_available": OCR_AVAILABLE}

@app.get("/health")
async def health_check():
    return {"status": "healthy", "timestamp": datetime.now().isoformat(), "ocr_available": OCR_AVAILABLE}

@app.post("/extract")
async def extract_pdf_endpoint(
    file: UploadFile = File(...),
    use_ocr: bool = Query(True, description="Enable OCR fallback for scanned pages"),
    force_ocr: bool = Query(False, description="Force OCR on all pages to extract text from images/logos"),
):
    if not file.filename.lower().endswith('.pdf'):
        raise HTTPException(status_code=400, detail="File must be a PDF")
    
    try:
        pdf_bytes = await file.read()
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Failed to read file: {str(e)}")
    
    if len(pdf_bytes) == 0:
        raise HTTPException(status_code=400, detail="Empty file uploaded")
    
    try:
        result = process_pdf(pdf_bytes, file.filename, use_ocr=use_ocr, force_ocr=force_ocr)
        return JSONResponse(content=result)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Processing error: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
