 #!/usr/bin/env python3
"""
simple_pdf_extract_test.py
--------------------------
Standalone PDF-to-JSON extraction script using PyMuPDF.
Tests extraction locally before running as API.

Usage:
    python simple_pdf_extract_test.py <input.pdf> <output.json>
    python simple_pdf_extract_test.py sample.pdf output.json --verbose
"""

import sys
import json
import argparse
from pathlib import Path
from datetime import datetime
from typing import Optional

try:
    import fitz  # PyMuPDF
except ImportError:
    print("ERROR: PyMuPDF not installed. Run: pip install PyMuPDF")
    sys.exit(1)

# Optional OCR support
try:
    import pytesseract
    from PIL import Image
    import io
    OCR_AVAILABLE = True
except ImportError:
    OCR_AVAILABLE = False


def extract_pdf_metadata(doc: fitz.Document) -> dict:
    """Extract PDF metadata (author, title, dates, etc.)"""
    meta = doc.metadata or {}
    return {
        "title": meta.get("title", ""),
        "author": meta.get("author", ""),
        "subject": meta.get("subject", ""),
        "keywords": meta.get("keywords", ""),
        "creator": meta.get("creator", ""),
        "producer": meta.get("producer", ""),
        "creation_date": meta.get("creationDate", ""),
        "modification_date": meta.get("modDate", ""),
        "page_count": doc.page_count,
        "is_encrypted": doc.is_encrypted,
        "is_pdf": doc.is_pdf,
    }


def extract_page_text(page: fitz.Page, use_ocr: bool = False, verbose: bool = False) -> dict:
    """
    Extract text from a single page.
    Falls back to OCR if text extraction yields nothing and OCR is enabled.
    """
    # Try native text extraction first
    text = page.get_text("text").strip()
    extraction_method = "native"
    
    # If no text found and OCR is available/enabled, try OCR
    if not text and use_ocr and OCR_AVAILABLE:
        if verbose:
            print(f"    Page {page.number + 1}: No native text, attempting OCR...")
        try:
            # Render page to image at 300 DPI for better OCR
            mat = fitz.Matrix(300/72, 300/72)
            pix = page.get_pixmap(matrix=mat)
            img_data = pix.tobytes("png")
            img = Image.open(io.BytesIO(img_data))
            text = pytesseract.image_to_string(img).strip()
            extraction_method = "ocr"
        except Exception as e:
            if verbose:
                print(f"    OCR failed: {e}")
            text = ""
            extraction_method = "ocr_failed"
    
    return {
        "page_number": page.number + 1,
        "text": text,
        "text_length": len(text),
        "extraction_method": extraction_method,
        "has_text": bool(text),
    }


def extract_page_images(page: fitz.Page, verbose: bool = False) -> list:
    """Extract information about images on a page."""
    images = []
    try:
        image_list = page.get_images(full=True)
        for img_index, img in enumerate(image_list):
            xref = img[0]
            images.append({
                "image_index": img_index + 1,
                "xref": xref,
                "width": img[2],
                "height": img[3],
                "bits_per_component": img[4],
                "colorspace": img[5],
            })
    except Exception as e:
        if verbose:
            print(f"    Error extracting images: {e}")
    return images


def extract_page_links(page: fitz.Page) -> list:
    """Extract hyperlinks from a page."""
    links = []
    try:
        for link in page.get_links():
            link_info = {
                "kind": link.get("kind", 0),
                "uri": link.get("uri", ""),
                "page": link.get("page", -1),
            }
            if link_info["uri"] or link_info["page"] >= 0:
                links.append(link_info)
    except Exception:
        pass
    return links


def extract_pdf(
    pdf_path: str,
    use_ocr: bool = True,
    verbose: bool = False
) -> dict:
    """
    Main extraction function.
    Returns structured JSON-compatible dict with all PDF content.
    """
    path = Path(pdf_path)
    if not path.exists():
        raise FileNotFoundError(f"PDF not found: {pdf_path}")
    
    if verbose:
        print(f"Opening: {pdf_path}")
        print(f"OCR enabled: {use_ocr and OCR_AVAILABLE}")
    
    doc = fitz.open(pdf_path)
    
    result = {
        "extraction_info": {
            "source_file": path.name,
            "source_path": str(path.absolute()),
            "extraction_timestamp": datetime.now().isoformat(),
            "pymupdf_version": fitz.version[0],
            "ocr_available": OCR_AVAILABLE,
            "ocr_used": use_ocr and OCR_AVAILABLE,
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
        if verbose:
            print(f"  Processing page {page_num + 1}/{doc.page_count}...")
        
        page = doc[page_num]
        
        page_data = extract_page_text(page, use_ocr=use_ocr, verbose=verbose)
        page_data["images"] = extract_page_images(page, verbose=verbose)
        page_data["links"] = extract_page_links(page)
        page_data["rotation"] = page.rotation
        page_data["width"] = page.rect.width
        page_data["height"] = page.rect.height
        
        result
