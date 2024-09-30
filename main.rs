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

#[cfg(test)]
mod git_test;

use args::Args;
use clap::Parser;
use std::io::{Error, ErrorKind};

fn main() -> Result<(), Error> {
    let args: Args = Args::parse();

    match args.task.as_str() {
        "create" => create_module_task(&args),
        "develop" => cmd_develop_modules::develop_modules_command(),
        _ => {
            eprintln!("Unknown task: {}", args.task);
            Err(Error::new(ErrorKind::InvalidInput, "Unknown task"))
        }
    }
}

fn create_module_task(args: &Args) -> Result<(), Error> {
    match &args.module {
        Some(module) => {
            let module_type = args.module_type.as_deref().unwrap_or("full");
            cmd_create_module::create_module(module, module_type)
        }
        None => {
            eprintln!("Module name is required for 'create' task");
            Err(Error::new(ErrorKind::InvalidInput, "Module name is required"))
        }
    }
}