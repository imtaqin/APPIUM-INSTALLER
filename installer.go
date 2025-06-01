package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Installer struct {
	installPath string
	projectPath string
	tempPath    string
}

func NewInstaller(installPath, projectPath string) *Installer {
	return &Installer{
		installPath: installPath,
		projectPath: projectPath,
		tempPath:    filepath.Join(os.TempDir(), "AppiumInstaller"),
	}
}

func (i *Installer) RunInstallation() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Setting up directories", i.setup},
		{"Installing Node.js", i.installNodeJS},
		{"Installing Java JDK", i.installJava},
		{"Installing ADB", i.checkAndInstallADB},
		{"Installing Appium", i.installAppium},
		{"Refreshing Environment", i.refreshEnvironment},
		{"Installing UiAutomator2", i.checkAndInstallUIAutomator2},
		{"Installing Dependencies", i.installDependencies},
		{"Setting Environment Variables", i.setEnvironmentVariables},
		{"Verifying Installation", i.verify},
		{"Cleaning up", i.cleanup},
	}

	totalSteps := len(steps)
	for idx, step := range steps {
		fmt.Printf("[%d/%d] %s...\n", idx+1, totalSteps, step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s failed: %v", step.name, err)
		}
		fmt.Printf("✅ %s completed\n", step.name)
		fmt.Println(strings.Repeat("-", 50))
	}

	return nil
}

func (i *Installer) setup() error {
	dirs := []string{i.installPath, i.tempPath}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		fmt.Printf("   📁 Created directory: %s\n", dir)
	}
	return nil
}

func (i *Installer) installNodeJS() error {
	if commandExists("node") {
		version, _ := runCommand("node", "--version")
		fmt.Printf("   ℹ️  Node.js already installed: %s\n", strings.TrimSpace(version))
		return nil
	}

	fmt.Println("   📥 Downloading Node.js...")
	nodeInstaller := filepath.Join(i.tempPath, "nodejs.msi")
	if err := downloadFile(NodeURL, nodeInstaller); err != nil {
		return fmt.Errorf("download failed: %v", err)
	}

	fmt.Println("   🔧 Installing Node.js (this may take a few minutes)...")
	if err := runCommandWait("msiexec.exe", "/i", nodeInstaller, "/quiet", "/norestart"); err != nil {
		return fmt.Errorf("installation failed: %v", err)
	}

	fmt.Printf("   ✅ Node.js installed successfully\n")
	return nil
}

func (i *Installer) installJava() error {
	javaPath := filepath.Join(i.installPath, "jdk-24.0.1")
	
	if _, err := os.Stat(javaPath); os.IsNotExist(err) {
		fmt.Println("   📥 Downloading Java JDK (this is a large file, please wait)...")
		javaZip := filepath.Join(i.tempPath, "openjdk-24.0.1.zip")
		if err := downloadFile(JavaURL, javaZip); err != nil {
			return fmt.Errorf("download failed: %v", err)
		}
		
		fmt.Println("   📦 Extracting Java JDK...")
		if err := extractZip(javaZip, i.installPath); err != nil {
			return fmt.Errorf("extraction failed: %v", err)
		}
		fmt.Printf("   ✅ Java JDK extracted to: %s\n", javaPath)
	} else {
		fmt.Printf("   ℹ️  Java JDK already exists at: %s\n", javaPath)
	}

	return nil
}

func (i *Installer) checkAndInstallADB() error {
	if commandExists("adb") {
		version, err := runCommand("adb", "version")
		if err == nil {
			fmt.Printf("   ℹ️  ADB already available: %s\n", strings.TrimSpace(strings.Split(version, "\n")[0]))
			return nil
		}
	}
	
	fmt.Println("   📥 Downloading Android Platform Tools...")
	androidPath := filepath.Join(i.installPath, "android-sdk")
	androidZip := filepath.Join(i.tempPath, "platform-tools.zip")
	if err := downloadFile(AndroidURL, androidZip); err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	
	if err := os.MkdirAll(androidPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	fmt.Println("   📦 Extracting Android Platform Tools...")
	if err := extractZip(androidZip, androidPath); err != nil {
		return fmt.Errorf("extraction failed: %v", err)
	}

	fmt.Printf("   ✅ ADB installed to: %s\n", filepath.Join(androidPath, "platform-tools"))
	return nil
}

func (i *Installer) installAppium() error {
	if commandExists("appium") {
		version, err := runCommand("appium", "--version")
		if err == nil {
			fmt.Printf("   ℹ️  Appium already installed: %s\n", strings.TrimSpace(version))
			return nil
		}
	}
	
	fmt.Println("   📦 Installing Appium globally using npm...")
	cmd := exec.Command("cmd", "/C", "npm install -g appium")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("   ❌ npm output: %s\n", string(output))
		return fmt.Errorf("npm install failed: %v", err)
	}

	fmt.Println("   ✅ Appium installed successfully")
	return nil
}

func (i *Installer) refreshEnvironment() error {
	fmt.Println("   🔄 Refreshing environment variables...")
	
	refreshEnvironmentCmd()
	
	npmPath, err := runCommand("cmd", "/C", "npm config get prefix")
	if err == nil {
		npmGlobalPath := strings.TrimSpace(npmPath)
		fmt.Printf("   📁 NPM Global Path: %s\n", npmGlobalPath)
		
		os.Setenv("PATH", os.Getenv("PATH")+";"+npmGlobalPath)
		if runtime.GOOS == "windows" {
			os.Setenv("PATH", os.Getenv("PATH")+";"+filepath.Join(npmGlobalPath, "node_modules", ".bin"))
		}
	}
	
	return nil
}

func (i *Installer) checkAndInstallUIAutomator2() error {
	fmt.Println("   🔍 Installing UiAutomator2 driver using npm...")
	
	cmd := exec.Command("cmd", "/C", "npm install -g appium-uiautomator2-driver")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("   ❌ npm output: %s\n", string(output))
		
		fmt.Println("   🔄 Trying alternative method with appium command...")
		appiumPaths := []string{
			"appium",
			"C:\\Users\\" + os.Getenv("USERNAME") + "\\AppData\\Roaming\\npm\\appium.cmd",
			"C:\\Program Files\\nodejs\\appium.cmd",
		}
		
		var appiumCmd string
		for _, path := range appiumPaths {
			if commandExists(path) || fileExists(path) {
				appiumCmd = path
				break
			}
		}
		
		if appiumCmd == "" {
			return fmt.Errorf("appium command not found, please restart terminal and run: appium driver install uiautomator2")
		}
		
		cmd2 := exec.Command("cmd", "/C", appiumCmd, "driver", "install", "uiautomator2")
		output2, err2 := cmd2.CombinedOutput()
		if err2 != nil {
			fmt.Printf("   ❌ appium output: %s\n", string(output2))
			return fmt.Errorf("UiAutomator2 installation failed, please restart terminal and run: appium driver install uiautomator2")
		}
	}

	fmt.Println("   ✅ UiAutomator2 driver installed successfully")
	return nil
}

func (i *Installer) installDependencies() error {
	if _, err := os.Stat(i.projectPath); os.IsNotExist(err) {
		fmt.Printf("   ℹ️  Project path %s does not exist, skipping npm install\n", i.projectPath)
		return nil
	}

	packageJsonPath := filepath.Join(i.projectPath, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		fmt.Println("   ℹ️  No package.json found, skipping npm install")
		return nil
	}

	fmt.Printf("   📦 Running npm install in %s...\n", i.projectPath)
	cmd := exec.Command("npm", "install")
	cmd.Dir = i.projectPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install failed: %v", err)
	}
	fmt.Println("   ✅ Dependencies installed successfully")
	return nil
}

func (i *Installer) setEnvironmentVariables() error {
	javaPath := filepath.Join(i.installPath, "jdk-24.0.1")
	androidPath := filepath.Join(i.installPath, "android-sdk")
	platformToolsPath := filepath.Join(androidPath, "platform-tools")

	envVars := map[string]string{
		"JAVA_HOME":       javaPath,
		"ANDROID_HOME":    androidPath,
		"ANDROID_SDK_ROOT": androidPath,
	}

	for name, value := range envVars {
		fmt.Printf("   🔧 Setting %s = %s\n", name, value)
		if err := setEnvVarPermanent(name, value); err != nil {
			fmt.Printf("   ⚠️  Warning: Failed to set %s: %v\n", name, err)
		}
	}

	pathsToAdd := []string{
		filepath.Join(javaPath, "bin"),
		platformToolsPath,
	}

	for _, path := range pathsToAdd {
		fmt.Printf("   🔧 Adding to PATH: %s\n", path)
		if err := addToPathPermanent(path); err != nil {
			fmt.Printf("   ⚠️  Warning: Failed to add to PATH: %v\n", err)
		}
	}

	return nil
}

func (i *Installer) verify() error {
	fmt.Println("   🔍 Verifying installations...")
	
	checks := map[string][]string{
		"Node.js": {"node", "--version"},
		"NPM":     {"npm", "--version"},
		"Java":    {"java", "--version"},
		"ADB":     {"adb", "version"},
	}

	allGood := true
	for name, cmd := range checks {
		if output, err := runCommand(cmd[0], cmd[1:]...); err == nil {
			version := strings.TrimSpace(strings.Split(output, "\n")[0])
			fmt.Printf("   ✅ %s: %s\n", name, version)
		} else {
			fmt.Printf("   ❌ %s: Not found in PATH (restart terminal required)\n", name)
			allGood = false
		}
	}

	appiumPaths := []string{
		"appium",
		"C:\\Users\\" + os.Getenv("USERNAME") + "\\AppData\\Roaming\\npm\\appium.cmd",
	}
	
	appiumFound := false
	for _, path := range appiumPaths {
		if output, err := runCommand(path, "--version"); err == nil {
			version := strings.TrimSpace(strings.Split(output, "\n")[0])
			fmt.Printf("   ✅ Appium: %s\n", version)
			appiumFound = true
			break
		}
	}
	
	if !appiumFound {
		fmt.Printf("   ⚠️  Appium: Not found in PATH (restart terminal required)\n")
		allGood = false
	}

	fmt.Println("\n   📁 Environment Variables:")
	envVars := []string{"JAVA_HOME", "ANDROID_HOME", "ANDROID_SDK_ROOT"}
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			fmt.Printf("   ✅ %s: %s\n", envVar, value)
		} else {
			defaultPath := filepath.Join(i.installPath, strings.ToLower(strings.Replace(envVar, "_", "-", -1)))
			if envVar == "JAVA_HOME" {
				defaultPath = filepath.Join(i.installPath, "jdk-24.0.1")
			} else {
				defaultPath = filepath.Join(i.installPath, "android-sdk")
			}
			fmt.Printf("   ⚠️  %s: Not set in current session (will be available after restart)\n", envVar)
			fmt.Printf("      Expected value: %s\n", defaultPath)
		}
	}

	if !allGood {
		fmt.Println("\n   ⚠️  Some tools may not be immediately available.")
		fmt.Println("      This is normal - restart your terminal to refresh environment variables.")
	}

	return nil
}

func (i *Installer) cleanup() error {
	fmt.Printf("   🧹 Removing temporary files from: %s\n", i.tempPath)
	if err := os.RemoveAll(i.tempPath); err != nil {
		fmt.Printf("   ⚠️  Warning: Failed to cleanup temp files: %v\n", err)
	}
	return nil
}