#[cfg(test)]
mod tests {
    use crate::git;
    use std::env;
    use std::path::{Path, PathBuf};

    #[test]
    fn test_find_git_root() {
        // Start from the current executable's directory and search upwards
        let mut current_dir = env::current_exe().unwrap().parent().unwrap().to_path_buf();
        
        while current_dir.parent().is_some() {
            if let Ok(git_root) = git::find_git_root_from_path(&current_dir) {
                // We found a git root
                assert!(Path::new(&git_root).is_absolute());
                assert!(Path::new(&git_root).join(".git").exists());
                return;
            }
            // Move up one directory
            current_dir = current_dir.parent().unwrap().to_path_buf();
        }
        
        panic!("No Git repository found in any parent directory");
    }

    #[test]
    fn test_find_git_root_not_in_repo() {
        // This test remains unchanged
        let original_dir = env::current_dir().unwrap();
        env::set_current_dir("/tmp").unwrap();

        let result = git::find_git_root();
        assert!(result.is_err());

        env::set_current_dir(original_dir).unwrap();
    }
}