package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("❌ This installer is designed for Windows only")
		waitForExit()
		os.Exit(1)
	}

	printHeader()
	
	installer := NewInstaller(getInstallPath(), getProjectPath())

	fmt.Printf("📁 Install Path: %s\n", installer.installPath)
	fmt.Printf("📂 Project Path: %s\n", installer.projectPath)
	fmt.Println()

	if err := installer.RunInstallation(); err != nil {
		fmt.Printf("\n❌ INSTALLATION FAILED!\n")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\n💡 Troubleshooting Tips:")
		fmt.Println("  • Make sure you're running as Administrator")
		fmt.Println("  • Check your internet connection")
		fmt.Println("  • Disable antivirus temporarily")
		fmt.Println("  • Try running the installer again")
		waitForExit()
		os.Exit(1)
	}

	printSuccessMessage()
	waitForExit()
}

func waitForExit() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Print("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func printHeader() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    APPIUM INSTALLER CLI                      ║")
	fmt.Println("║                   Windows Edition                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printSuccessMessage() {
	fmt.Println("\n🎉 INSTALLATION COMPLETED SUCCESSFULLY!")
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    INSTALLATION SUMMARY                     ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Println("║ ✅ Node.js & NPM installed                                  ║")
	fmt.Println("║ ✅ Java JDK installed                                       ║")
	fmt.Println("║ ✅ Android ADB installed                                    ║")
	fmt.Println("║ ✅ Appium installed globally                                ║")
	fmt.Println("║ ✅ UiAutomator2 driver installed                           ║")
	fmt.Println("║ ✅ Environment variables configured                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	fmt.Println("\n📋 NEXT STEPS:")
	fmt.Println("  1. 🔄 Restart your terminal (IMPORTANT!)")
	fmt.Println("  2. 🚀 Start Appium server: appium")
	fmt.Println("  3. 🔌 Connect device/emulator: adb devices")
	fmt.Println("  4. 📱 Navigate to project and run your tests")
	
	fmt.Println("\n🔧 VERIFICATION COMMANDS:")
	fmt.Println("  • Check Java: java -version")
	fmt.Println("  • Check ADB: adb version")
	fmt.Println("  • Check Appium: appium --version")
	fmt.Println("  • Check Environment: echo %JAVA_HOME%")
	
	fmt.Println("\n🚀 Happy Mobile Testing!")
}

func getInstallPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return DefaultInstallPath
}

func getProjectPath() string {
	if len(os.Args) > 2 {
		return os.Args[2]
	}
	return "C:\\AppiumProject"
}