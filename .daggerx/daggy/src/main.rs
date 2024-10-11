#[cfg(test)]
mod git_test;

mod args;
mod command_utils;
mod configuration;
mod dagger_commands;
mod dagger_json;
mod dagger_utils;
mod git;
mod github_actions;
mod naming;
mod readme_and_docs;
mod templating;
mod utils;
mod cmd_create_module;
mod cmd_develop_modules;

use args::Args;
use clap::Parser;
use std::io::{Error, ErrorKind};

const SUPPORTED_MODULE_TYPES: [&str; 2] = ["full", "light"];

fn main() -> Result<(), Error> {
    let args: Args = Args::parse();

    match args.task.as_str() {
        "create" => create_module_task(&args),
        "sync" => sync_modules_task(),
        "inspect" => inspect_modules_task(),
        "develop" => cmd_develop_modules::develop_modules(),
        _ => {
            eprintln!("Unknown task: {}", args.task);
            Err(Error::new(ErrorKind::InvalidInput, "Unknown task"))
        }
    }
}

fn sync_modules_task() -> Result<(), Error> {
    println!("Syncing modules...");
    Ok(())
}

fn inspect_modules_task() -> Result<(), Error> {
    println!("Inspecting modules...");
    Ok(())
}

fn create_module_task(args: &Args) -> Result<(), Error> {
    match &args.module {
        Some(module) => {
            cmd_create_module::create_module(module, args.module_type.as_deref().unwrap_or("full"))
        }
        None => {
            eprintln!("Module name is required for 'create' task");
            Err(Error::new(ErrorKind::InvalidInput, "Module name is required"))
        }
    }
}