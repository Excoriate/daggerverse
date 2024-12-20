#[cfg(test)]
mod git_test;

mod args;
mod cmd_create_module;
mod cmd_develop_modules;
mod command_utils;
mod configuration;
mod dagger_commands;
mod dagger_json;
mod dagger_utils;
mod git;
mod github_actions;
mod naming;
mod readme_and_docs;
mod sync_inspect;
mod templating;
mod utils;

use args::Args;
use clap::Parser;
use std::io::{Error, ErrorKind};

const SUPPORTED_MODULE_TYPES: [&str; 2] = ["full", "light"];

fn main() -> Result<(), Error> {
    let args: Args = Args::parse();

    match args.task.as_str() {
        "create" => create_module_task(&args),
        "sync" => sync_inspect::sync_modules_task(&args),
        "inspect" => sync_inspect::inspect_modules_task(&args),
        "develop" => cmd_develop_modules::develop_modules(),
        _ => {
            eprintln!("Unknown task: {}", args.task);
            Err(Error::new(ErrorKind::InvalidInput, "Unknown task"))
        }
    }
}

fn create_module_task(args: &Args) -> Result<(), Error> {
    match &args.module {
        Some(module) => {
            cmd_create_module::create_module(module, args.module_type.as_deref().unwrap_or("full"))
        }
        None => {
            eprintln!("Module name is required for 'create' task");
            Err(Error::new(
                ErrorKind::InvalidInput,
                "Module name is required",
            ))
        }
    }
}
