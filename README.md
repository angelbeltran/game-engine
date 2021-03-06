# game-engine
Pet project. Defining game rules and dynamics via Protobufs in Go.

# Roadmap
- [x] reference input in rules
- [x] multiple effects (list)
- [x] define responses (not just errors)
- [x] variadic function (eg, "and" and "or")
- [ ] dynamic state/input field references
- [ ] fail effects if input or state field is not initialized
- [ ] support binding errors per rule at anywhere in the rule tree
- [ ] initialize state per game
    - [ ] require a new special 'new game' action
        - [ ] initial new 'fully initialized' state on request (eg, no nil/null)
    - [ ] manage game sessions
- [ ] plugin support
    - [ ] db
        - [ ] sql plugin
        - [ ] mongodb?
    - [ ] events
        - [ ] publications
            - [ ] state changes
        - [ ] subscriptions
            - [ ] arbitrary external sources
- [ ] full type support
    - [ ] enums
    - [ ] maps
    - [ ] oneofs
    - [ ] various numerical types
    - [ ] bytes
    - [ ] others?

## Dev niceties
- [x] replace as many template functions with templates.

## Examples

Check out the example directory for an idea of how to use the tool and what it's supposed to do.
