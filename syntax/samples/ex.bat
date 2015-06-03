    @ECHO OFF

    VER | FIND "Windows 2000" >NUL

    IF ERRORLEVEL 1 GOTO Syntax

     

    SETLOCAL ENABLEEXTENSIONS

     

    :: Parameter check

    IF [%2]==[] (

    	ENDLOCAL

    	GOTO Syntax

    )

     

    SET Group=

    SET Domain=

    NET GROUP %1 /DOMAIN >NUL 2>NUL

    IF NOT ERRORLEVEL 1 (

    	SET Group=GROUP

    	SET Domain=/DOMAIN

    )

    NET GROUP %1 >NUL 2>NUL

    IF NOT ERRORLEVEL 1 (

    	SET Group=GROUP

    	SET Domain=

    )

    NET LOCALGROUP %1 /DOMAIN >NUL 2>NUL

    IF NOT ERRORLEVEL 1 (

    	SET Group=LOCALGROUP

    	SET Domain=/DOMAIN

    )

    NET LOCALGROUP %1 >NUL 2>NUL

    IF NOT ERRORLEVEL 1 (

    	SET Group=LOCALGROUP

    	SET Domain=

    )

    IF NOT DEFINED Group (

    	ECHO.

    	ECHO "%~1" is not a valid NT group

    	ENDLOCAL

    	GOTO Syntax

    )

     

    :: Main command

    FOR /F "skip=8 tokens=* delims=" %%A IN ('NET %Group% %1 %Domain% ^| FIND /V "The command"') DO FOR /F "tokens=1-3 delims=	 " %%a IN ('ECHO.%%A') DO (

    	CALL :Command "%%a" %*

    	IF NOT "%%b"=="" CALL :Command "%%b" %*

    	IF NOT "%%c"=="" CALL :Command "%%c" %*

    )

    ENDLOCAL

    GOTO End

     

     

    :Command

    SETLOCAL

    SET User$=%~1

    SHIFT

    :Loop

    SHIFT

    IF [%1]==[] GOTO Continue

    IF [%1]==[ ] GOTO Continue

    IF [%~1]==[#] (SET command$=%command$% %User$%) ELSE (SET command$=%command$% %1)

    GOTO Loop

    :Continue

    IF "%command$%"==" " GOTO Syntax

    CALL %command$%

    ENDLOCAL

    GOTO:EOF

     

     

    :Syntax

    ECHO.

    ECHO 4AllMembers,  Version 2.11 for Windows 2000

    ECHO Execute a command once for each member of a local, domain local or global group

    ECHO Written by Rob van der Woude

    ECHO http://www.robvanderwoude.com

    ECHO.

    IF "%OS%"=="Windows_NT" ECHO Usage:  %~n0  ^<group_name^>  ^<any_command^>  [ ^<parameters^> ]

    IF NOT "%OS%"=="Windows_NT" ECHO Usage:  4AllMembers  {group_name}  {my_command}  [ {parameters} ]

    ECHO.

    ECHO Any_command will be executed once for each member of group_name

    ECHO Command line parameter(s) "#" will be substituted by user ID

     

    :End

     
