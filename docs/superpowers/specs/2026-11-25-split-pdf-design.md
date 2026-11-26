# Split PDF into pages — design

**Date:** 2026-11-25 · **Status:** approved by owner · **Size:** small change

- **Why:** the owner wants every page of a PDF as its own PDF file.
- **Menu:** right-click a PDF → **Split into pages**, next to Reduce size.
- **Output:** `report-1.pdf`, `report-2.pdf`, … next to the original, the same
  rule as PDF → PNG. If any name is taken, the whole set uses the next free
  number: `report (1)-1.pdf`, …
- **How:** `mutool merge` copies each page as it is, with no re-compression.
  The page count comes from `mutool show`, which also refuses a
  password-protected PDF, so the user gets the usual clear message.
- **Messages:**
  - one-page PDF: "This PDF has only one page" (no file made);
  - Done notification: "Done: report-1.pdf + 11 more (12 pages)".
- **Alternative not chosen:** a `report (pages)\` folder. Tidier for big
  PDFs, but different from how PDF → PNG works.
