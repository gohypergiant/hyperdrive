import traceback
from fastapi import APIRouter, File, UploadFile
from fastapp.services.hyperpackage import hyperpackage_upload_path, hyperpackage_upload_file, extract_hyperpackage, clear_uploads

router = APIRouter()


@router.post("/upload-hyperpack")
def upload(file: UploadFile = File(...)):
    clear_uploads()
    try:
        with open(hyperpackage_upload_file, 'wb') as f:
            while contents := file.file.read(1024 * 1024):
                f.write(contents)
        extract_hyperpackage(hyperpackage_upload_path, hyperpackage_upload_file)

    except Exception as e:
        print(e)
        return {"message": "There was an error uploading the file", "success": False, "error": str(e)}
    finally:
        file.file.close()

    return {"message": f"Successfully uploaded {file.filename}", "success": True}
