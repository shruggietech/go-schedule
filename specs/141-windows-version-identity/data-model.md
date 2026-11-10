# Data model

Selected version has original label, numeric three-part value and development status. Release components obey MSI limits (255, 255, 65535). Development maps to numeric 0.0.1 and explicit comment. Candidate metadata joins executable fixed/string resources, embedded Go linker version, MSI ProductVersion and its GUI File.Version row. No persistent product state changes.
