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

use std::borrow::Borrow;
use std::collections::hash_map::{Entry, Iter, Keys, Values};
use std::hash::Hash;

/// Threadsafe variant of std::collections::Vector<T>
pub trait List<T> {
    fn new() -> Self;
    fn with_capacity(capacity: usize) -> Self;
    fn push(&mut self, value: T);
    fn pop(&mut self) -> T;
    fn len(&self) -> usize;
    fn capacity(&self) -> usize;
    fn is_empty(&self) -> bool;
    fn clear(&mut self);
    fn remove(&mut self, index: usize) -> T;
}

/// Threadsafe variant of std::collections::hash_map::HashMap<K, V>
// todo (sienna): implement crossbeam.SkipMap & BTreeMap
pub trait Map<K, V> {
    fn new() -> Self;
    fn with_capacity(capacity: usize) -> Self;

    fn contains_key<Q>(&self, k: &Q) -> bool
    where
        K: Borrow<Q>,
        Q: Hash + Eq + ?Sized;

    fn entry(&mut self, key: K) -> Entry<'_, K, V>;
    fn get<Q>(&self, k: &Q) -> Option<&V>
    where
        K: Borrow<Q>,
        Q: Hash + Eq + ?Sized;

    fn insert(&mut self, k: K, v: V) -> Option<V>;
    fn iter(&self) -> Iter<'_, K, V>;
    fn remove<Q>(&mut self, k: &Q) -> Option<V>
    where
        K: Borrow<Q>,
        Q: Hash + Eq + ?Sized;

    fn keys(&self) -> Keys<'_, K, V>;
    fn values(&self) -> Values<'_, K, V>;
}
