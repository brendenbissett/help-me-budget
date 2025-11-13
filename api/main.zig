const std = @import("std");

pub fn main() void {
    const appName = "Help-Me-Budget";
    const version = "0.0.1";

    std.debug.print("Welcome to {s} API (Version: {s})", .{ appName, version });
    return;
}
