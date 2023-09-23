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

use std::ffi::OsString;
use clap::Parser;

#[derive(Parser, Debug)]
#[command(name = "server")]
pub struct ServerConfig {
    /// Port to run this node on. It must be the same for all nodes in the cluster
    #[arg(short, long, default_value = "8080")]
    port: Option<i32>,

    /// The root directory where Pleiades will store all of it's data.
    #[arg(long, default_value = "/var/pleiades")]
    data_dir: Option<OsString>,
}
