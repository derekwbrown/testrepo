call gradlew.bat build
if not "%ERRORLEVEL%" == "0" goto :EOF
java -jar signapk.jar platform.x509.pem platform.pk8 app\build\outputs\apk\release\app-release-unsigned.apk app-release-signed.apk
if not "%ERRORLEVEL%" == "0" goto :EOF
adb install -r app-release-signed.apk