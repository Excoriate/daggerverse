use std::fs;
use std::io::Error;
use std::path::Path;

pub fn copy_dir_all(src: impl AsRef<Path>, dst: impl AsRef<Path>) -> Result<(), Error> {
    fs::create_dir_all(&dst)?;
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let ty = entry.file_type()?;
        if ty.is_dir() {
            copy_dir_all(entry.path(), dst.as_ref().join(entry.file_name()))?;
        } else {
            fs::copy(entry.path(), dst.as_ref().join(entry.file_name()))?;
        }
    }
    Ok(())
}

pub fn copy_dir_recursive(src: &Path, dest: &Path, module_cfg: &crate::configuration::NewDaggerModule) -> Result<(), Error> {
    if !dest.exists() {
        fs::create_dir_all(dest)?;
    }

    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let file_type = entry.file_type()?;
        let src_path = entry.path();
        let mut file_name = entry.file_name().to_string_lossy().to_string();

        // Remove .tmpl extension if present
        if file_name.ends_with(".tmpl") {
            file_name = file_name.trim_end_matches(".tmpl").to_string();
        }

        let dest_path = dest.join(&file_name);

        if file_type.is_dir() {
            copy_dir_recursive(&src_path, &dest_path, module_cfg)?;
        } else {
            let content = fs::read_to_string(&src_path)?;
            let processed_content = crate::templating::process_template_content(&content, module_cfg);
            fs::write(dest_path, processed_content)?;
        }
    }

    Ok(())
}
