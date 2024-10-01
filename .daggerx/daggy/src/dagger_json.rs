use crate::configuration::NewDaggerModule;
use serde_json::Value;
use std::fs;
use std::io::{Error, ErrorKind};

const TEST_JSON_CONTENT: &str = r#"
{
    "exclude": [
        "../../.direnv",
        "../../.devenv",
        "../../.vscode",
        "../../.idea",
        "../../.trunk",
        "../../go.work",
        "../../go.work.sum"
    ]
}
"#;

const EXAMPLES_JSON_CONTENT: &str = r#"
{
    "exclude": [
        "../../.direnv",
        "../../.devenv",
        "../../.vscode",
        "../../.idea",
        "../../.trunk",
        "../../go.work",
        "../../go.work.sum"
    ]
}
"#;

const DAGGER_JSON_CONTENT: &str = r#"
{
    "exclude": [
        "../.direnv",
        "../.devenv",
        "../.vscode",
        "../.idea",
        "../.trunk",
        "../go.work",
        "../go.work.sum"
    ]
}
"#;

pub fn update_tests_dagger_json(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let dagger_json_path = format!("{}/tests/dagger.json", module_cfg.path);
    let mut json_content: Value = fs::read_to_string(&dagger_json_path)
        .map_err(|e| {
            Error::new(
                ErrorKind::Other,
                format!("Failed to read tests/dagger.json: {}", e),
            )
        })
        .and_then(|content| {
            serde_json::from_str(&content).map_err(|e| {
                Error::new(
                    ErrorKind::Other,
                    format!("Failed to parse tests/dagger.json: {}", e),
                )
            })
        })?;

    let test_json_content: Value = serde_json::from_str(TEST_JSON_CONTENT).map_err(|e| {
        Error::new(
            ErrorKind::Other,
            format!("Failed to parse test JSON content: {}", e),
        )
    })?;

    json_content["exclude"] = test_json_content["exclude"].clone();

    fs::write(
        dagger_json_path,
        serde_json::to_string_pretty(&json_content)?,
    )
    .map_err(|e| {
        Error::new(
            ErrorKind::Other,
            format!("Failed to write updated tests/dagger.json: {}", e),
        )
    })?;

    Ok(())
}

pub fn update_examples_dagger_json(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let dagger_json_path = format!("{}/examples/go/dagger.json", module_cfg.path);
    let mut json_content: Value = fs::read_to_string(&dagger_json_path)
        .map_err(|e| {
            Error::new(
                ErrorKind::Other,
                format!("Failed to read examples/go/dagger.json: {}", e),
            )
        })
        .and_then(|content| {
            serde_json::from_str(&content).map_err(|e| {
                Error::new(
                    ErrorKind::Other,
                    format!("Failed to parse examples/go/dagger.json: {}", e),
                )
            })
        })?;

    let examples_json_content: Value =
        serde_json::from_str(EXAMPLES_JSON_CONTENT).map_err(|e| {
            Error::new(
                ErrorKind::Other,
                format!("Failed to parse examples JSON content: {}", e),
            )
        })?;

    json_content["exclude"] = examples_json_content["exclude"].clone();

    fs::write(
        dagger_json_path,
        serde_json::to_string_pretty(&json_content)?,
    )
    .map_err(|e| {
        Error::new(
            ErrorKind::Other,
            format!("Failed to write updated examples/go/dagger.json: {}", e),
        )
    })?;

    Ok(())
}

pub fn update_dagger_json(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let dagger_json_path = format!("{}/dagger.json", module_cfg.path);

    let mut json_content: Value = fs::read_to_string(&dagger_json_path)
        .map_err(|e| {
            Error::new(
                ErrorKind::Other,
                format!("Failed to read dagger.json: {}", e),
            )
        })
        .and_then(|content| {
            serde_json::from_str(&content).map_err(|e| {
                Error::new(
                    ErrorKind::Other,
                    format!("Failed to parse dagger.json: {}", e),
                )
            })
        })?;

    let dagger_json_content: Value = serde_json::from_str(DAGGER_JSON_CONTENT).map_err(|e| {
        Error::new(
            ErrorKind::Other,
            format!("Failed to parse dagger JSON content: {}", e),
        )
    })?;

    json_content["exclude"] = dagger_json_content["exclude"].clone();

    fs::write(
        dagger_json_path,
        serde_json::to_string_pretty(&json_content)?,
    )
    .map_err(|e| {
        Error::new(
            ErrorKind::Other,
            format!("Failed to write updated dagger.json: {}", e),
        )
    })?;

    Ok(())
}
