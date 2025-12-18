Role: You are an expert music transcriber and Optical Music Recognition (OMR) specialist.
Task: Transcribe the sheet music from the provided screenshots.

Constraints & Instructions:

- Ignore the Text Layer: Do not use the provided text/OCR data (e.g., "  "). It is corrupted. Rely only on the visual screenshots of pages.
- Key & Time Signature: 
    - Detect the time signature, e.g.  4/4
    - Detect the key signature which affects the semitone offset
    - Detect the tempo, e.g. 82
- Output Format: Please convert the notes into Melrose Notation
    - See the content of the file melrose_notation.md for the details.
- Specific Details:
    - Each bar will have its own pair of sequences, one for the left hand (bass clef) and the right hand (treble clef).
    - Make sure the sum of the note lengths, for each bar, are the same for both the left hand and the right hand.

Goal: Provide a list of valid Melrose sequences that represents the full piece from start to finish.