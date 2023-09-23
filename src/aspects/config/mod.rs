/*
 *   Pleiades Source Code
 *   Copyright (C) 2023 Sienna Lloyd
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>
 */

use clap::{Parser, Subcommand, ArgAction};

use crate::aspects::config::server::ServerConfig;

pub mod server;

#[derive(Parser, Debug)]
#[command(version, about)]
pub struct Root {
    /// Skip TLS checks on API calls.
    #[arg(
    short,
    long,
    env,
    global = true,
    name = "skip-tls-verify",
    )]
    pub skip_tls_verify: bool,

    /// Enable debug logging
    #[arg(short, long, global = true, conflicts_with = "trace")]
    pub debug: bool,

    /// Enable trace logging
    #[arg(short, long, global = true, conflicts_with = "debug")]
    pub trace: bool,

    #[command(subcommand)]
    command: Option<Commands>
}

#[derive(Subcommand, Debug)]
pub enum Commands {
    /// Run a local server node
    Server(ServerConfig)
}