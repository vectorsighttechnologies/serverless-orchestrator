import os
import zipfile
import subprocess

def pack():
    # 1. Build the binary
    print("Building Go Lambda binary (linux/arm64)...")
    env = os.environ.copy()
    env["GOOS"] = "linux"
    env["GOARCH"] = "arm64"
    
    # Ensure build dir exists
    os.makedirs(".build", exist_ok=True)
    
    subprocess.run([
        "go", "build", 
        "-ldflags=-s -w", 
        "-o", ".build/bootstrap", 
        "./cmd/lambda/"
    ], env=env, check=True)

    # 2. Package into function.zip with 0755 permissions
    print("Packaging bootstrap into function.zip...")
    func_zip_path = "function.zip"
    if os.path.exists(func_zip_path):
        os.remove(func_zip_path)

    with zipfile.ZipFile(func_zip_path, "w", zipfile.ZIP_DEFLATED) as z:
        info = zipfile.ZipInfo("bootstrap")
        info.create_system = 3 # Unix
        info.external_attr = 0o100755 << 16
        with open(".build/bootstrap", "rb") as f:
            z.writestr(info, f.read())
            
    print("Successfully created function.zip")

    # 3. Package into orchestrator-sam.zip in frontend/public
    print("Packaging template.yaml and function.zip into orchestrator-sam.zip...")
    sam_zip_dir = "../frontend/public"
    os.makedirs(sam_zip_dir, exist_ok=True)
    sam_zip_path = os.path.join(sam_zip_dir, "orchestrator-sam.zip")
    
    if os.path.exists(sam_zip_path):
        os.remove(sam_zip_path)
        
    with zipfile.ZipFile(sam_zip_path, "w", zipfile.ZIP_DEFLATED) as z:
        # Add template.yaml
        z.write("template.yaml", "template.yaml")
        # Add function.zip
        z.write("function.zip", "function.zip")
        
    print(f"Successfully created {sam_zip_path}")

if __name__ == "__main__":
    pack()
