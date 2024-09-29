use std::io::{Error, ErrorKind};
use std::process::{Command, Stdio};

pub fn run_dagger_develop(dir: &str) -> Result<(), Error> {
    let output = Command::new("dagger")
        .arg("develop")
        .current_dir(dir)
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()?;

    if !output.status.success() {
        return Err(Error::new(
            ErrorKind::Other,
            format!("dagger develop failed in directory: {}", dir),
        ));
    }

    Ok(())
}
