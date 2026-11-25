@echo off
rem Installs QuickConvert for the current user (no admin needed).
setlocal
set "DEST=%LOCALAPPDATA%\QuickConvert"

rem 0. Stop if run from inside the zip (Windows shows the zip like a folder).
if not exist "%~dp0quickconvert.exe" goto :notunzipped

rem 1. Remove menu entries from any older version.
reg delete "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.extract" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.clean" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.shrink" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.split" /f >nul 2>&1
reg delete "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.extract" /f >nul 2>&1

rem 2. Copy the program files (skipped when run from the install folder).
if /I not "%~dp0"=="%DEST%\" robocopy "%~dp0." "%DEST%" /E /NFL /NDL /NJH /NJS /NP >nul
powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-ChildItem -LiteralPath $env:DEST -Recurse | Unblock-File"

rem 3. Add the right-click menu entries.
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.3gp\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\05ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\05ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert\shell\05ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aac\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\05ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\05ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert\shell\05ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aif\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\05ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\05ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert\shell\05ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.aiff\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\05ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\05ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.amr\shell\QuickConvert\shell\05ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avi\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\03webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\03webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\03webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\04pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\04pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\04pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\05ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\05ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.avif\shell\QuickConvert\shell\05ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bat\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\03webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\03webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\03webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\04pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\04pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\04pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\05ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\05ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.bmp\shell\QuickConvert\shell\05ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.c\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.cpp\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.css\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert\shell\01pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert\shell\01pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert\shell\01pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert\shell\02md" /v MUIVerb /d "MD" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert\shell\02md" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.epub\shell\QuickConvert\shell\02md\command" /ve /d "\"%DEST%\quickconvert.exe\" md \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\04ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\04ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert\shell\04ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flac\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.flv\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\03webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\03webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\03webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\04pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\04pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\04pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\05ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\05ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\05ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\06mp4" /v MUIVerb /d "MP4 (video)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\06mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.gif\shell\QuickConvert\shell\06mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.h\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\03webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\03webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\03webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\04pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\04pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\04pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\05ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\05ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert\shell\05ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.heic\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.html\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.java\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\02webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\02webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\02webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\03pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\03pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\03pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\04ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\04ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert\shell\04ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpeg\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\02webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\02webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\02webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\03pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\03pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\03pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\04ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\04ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert\shell\04ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.jpg\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.js\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.json\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m2ts\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\02wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\02wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\02wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\03flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\03flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\03flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\04ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\04ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert\shell\04ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4a\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.m4v\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.md\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mkv\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\01m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\01m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\01m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\02wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\02wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\02wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\03flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\03flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\03flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\04ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\04ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert\shell\04ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp3\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\01whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\01whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\01whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\02mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\02mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\02mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\03m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\03m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\03m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\04gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\04gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert\shell\04gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mp4\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mpg\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.mts\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ogg\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\05ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\05ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert\shell\05ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.opus\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\03md" /v MUIVerb /d "MD" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\03md" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert\shell\03md\command" /ve /d "\"%DEST%\quickconvert.exe\" md \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.shrink" /v MUIVerb /d "Reduce size" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.shrink" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.shrink\command" /ve /d "\"%DEST%\quickconvert.exe\" shrink \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.split" /v MUIVerb /d "Split into pages" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.split" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.split\command" /ve /d "\"%DEST%\quickconvert.exe\" split \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.extract" /v MUIVerb /d "Extract pages..." /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.extract" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\QuickConvert.extract\command" /ve /d "\"%DEST%\quickconvert.exe\" extract \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\01jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\01jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\01jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\02webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\02webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\02webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\03pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\03pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\03pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\04ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\04ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert\shell\04ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.png\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.py\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.sh\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.svg\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\03webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\03webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\03webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\04pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\04pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\04pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\05ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\05ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tif\shell\QuickConvert\shell\05ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\03webp" /v MUIVerb /d "WEBP" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\03webp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\03webp\command" /ve /d "\"%DEST%\quickconvert.exe\" webp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\04pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\04pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\04pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\05ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\05ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.tiff\shell\QuickConvert\shell\05ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert\shell\01txt" /v MUIVerb /d ".txt" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert\shell\01txt" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.ts\shell\QuickConvert\shell\01txt\command" /ve /d "\"%DEST%\quickconvert.exe\" txt \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\01code-py" /v MUIVerb /d ".py" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\01code-py" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\01code-py\command" /ve /d "\"%DEST%\quickconvert.exe\" code-py \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\02code-m" /v MUIVerb /d ".m" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\02code-m" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\02code-m\command" /ve /d "\"%DEST%\quickconvert.exe\" code-m \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\03code-c" /v MUIVerb /d ".c" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\03code-c" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\03code-c\command" /ve /d "\"%DEST%\quickconvert.exe\" code-c \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\04code-cpp" /v MUIVerb /d ".cpp" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\04code-cpp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\04code-cpp\command" /ve /d "\"%DEST%\quickconvert.exe\" code-cpp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\05code-h" /v MUIVerb /d ".h" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\05code-h" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\05code-h\command" /ve /d "\"%DEST%\quickconvert.exe\" code-h \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\06code-js" /v MUIVerb /d ".js" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\06code-js" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\06code-js\command" /ve /d "\"%DEST%\quickconvert.exe\" code-js \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\07code-ts" /v MUIVerb /d ".ts" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\07code-ts" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\07code-ts\command" /ve /d "\"%DEST%\quickconvert.exe\" code-ts \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\08code-html" /v MUIVerb /d ".html" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\08code-html" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\08code-html\command" /ve /d "\"%DEST%\quickconvert.exe\" code-html \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\09code-css" /v MUIVerb /d ".css" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\09code-css" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\09code-css\command" /ve /d "\"%DEST%\quickconvert.exe\" code-css \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\10code-java" /v MUIVerb /d ".java" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\10code-java" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\10code-java\command" /ve /d "\"%DEST%\quickconvert.exe\" code-java \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\11code-json" /v MUIVerb /d ".json" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\11code-json" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\11code-json\command" /ve /d "\"%DEST%\quickconvert.exe\" code-json \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\12code-md" /v MUIVerb /d ".md" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\12code-md" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\12code-md\command" /ve /d "\"%DEST%\quickconvert.exe\" code-md \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\13code-bat" /v MUIVerb /d ".bat" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\13code-bat" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\13code-bat\command" /ve /d "\"%DEST%\quickconvert.exe\" code-bat \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\14code-sh" /v MUIVerb /d ".sh" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\14code-sh" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.txt\shell\QuickConvert\shell\14code-sh\command" /ve /d "\"%DEST%\quickconvert.exe\" code-sh \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\03flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\03flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\03flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\04ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\04ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert\shell\04ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wav\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webm\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\01png" /v MUIVerb /d "PNG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\01png" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\01png\command" /ve /d "\"%DEST%\quickconvert.exe\" png \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\02jpg" /v MUIVerb /d "JPG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\02jpg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\02jpg\command" /ve /d "\"%DEST%\quickconvert.exe\" jpg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\03pdf" /v MUIVerb /d "PDF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\03pdf" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\03pdf\command" /ve /d "\"%DEST%\quickconvert.exe\" pdf \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\04ico" /v MUIVerb /d "ICO (icon)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\04ico" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert\shell\04ico\command" /ve /d "\"%DEST%\quickconvert.exe\" ico \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.webp\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\01mp3" /v MUIVerb /d "MP3" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\01mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\01mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\02m4a" /v MUIVerb /d "M4A" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\02m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\02m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\03wav" /v MUIVerb /d "WAV" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\03wav" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\03wav\command" /ve /d "\"%DEST%\quickconvert.exe\" wav \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\04flac" /v MUIVerb /d "FLAC" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\04flac" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\04flac\command" /ve /d "\"%DEST%\quickconvert.exe\" flac \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\05ogg" /v MUIVerb /d "OGG" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\05ogg" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert\shell\05ogg\command" /ve /d "\"%DEST%\quickconvert.exe\" ogg \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wma\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert" /v MUIVerb /d "Convert to" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert" /v SubCommands /d "" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\01mp4" /v MUIVerb /d "MP4 (full quality, for editing)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\01mp4" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\02whatsapp" /v MUIVerb /d "MP4 (WhatsApp / any device)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\02whatsapp" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\02whatsapp\command" /ve /d "\"%DEST%\quickconvert.exe\" whatsapp \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\03mp3" /v MUIVerb /d "MP3 (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\03mp3" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\03mp3\command" /ve /d "\"%DEST%\quickconvert.exe\" mp3 \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\04m4a" /v MUIVerb /d "M4A (audio only)" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\04m4a" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\04m4a\command" /ve /d "\"%DEST%\quickconvert.exe\" m4a \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\05gif" /v MUIVerb /d "GIF" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\05gif" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert\shell\05gif\command" /ve /d "\"%DEST%\quickconvert.exe\" gif \"%%1\"" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.clean" /v MultiSelectModel /d "Player" /f >nul
reg add "HKCU\Software\Classes\SystemFileAssociations\.wmv\shell\QuickConvert.clean\command" /ve /d "\"%DEST%\quickconvert.exe\" clean \"%%1\"" /f >nul

rem 4. Add the "Send to" shortcuts (cmd cannot make shortcuts, PowerShell can).
powershell -NoProfile -ExecutionPolicy Bypass -Command "$s=(New-Object -ComObject WScript.Shell).CreateShortcut([Environment]::GetFolderPath('SendTo')+'\Images to PDF.lnk'); $s.TargetPath=$env:DEST+'\quickconvert.exe'; $s.Arguments='imgpdf'; $s.Save()"
powershell -NoProfile -ExecutionPolicy Bypass -Command "$s=(New-Object -ComObject WScript.Shell).CreateShortcut([Environment]::GetFolderPath('SendTo')+'\Merge PDFs.lnk'); $s.TargetPath=$env:DEST+'\quickconvert.exe'; $s.Arguments='merge'; $s.Save()"

echo QuickConvert installed.
echo Right-click a file ^> Show more options ^> Convert to.
pause
exit /b 0

:notunzipped
echo Please unzip QuickConvert first, then run install.cmd from the unzipped folder.
pause
exit /b 1
