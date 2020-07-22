# game-engine
Pet project. Defining game rules and dynamics via Protobufs in Go.

# Roadmap
- [ ] reference input in rules
- [x] multiple effects (list)
- [ ] define responses (not just errors)
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
