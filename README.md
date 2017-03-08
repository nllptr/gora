### GÖRA

[ˈʝɶːˌra]

Verb
1. do
2. make
3. manufacture, produce

"Göra" means "do" in Swedish. Gora is a Golang library for parsing and writing TODO-lists in a human-readable format. An example of such a list can be seen below.

    My evil plan
    ============
    - [ ] Build doomsday weapon
      - [X] Buy stuff
      - [ ] Built it
    - [ ] Demand ransom
      - [ ] One million dollars?
    - [ ] Avoid revealing plan to secret agent
    - [ ] Win

## Usage
Assuming you have the list above in a byte slice, parsing a list is as easy as

    l := gora.Parse(bytes)

To create a new list, add items, and generate its textual representation, do the following

    l := gora.New("My list")
    l.Add("New item 1")

Adding an item to a list or list item returns the newly created item:

    li := l.Add("New item 2")

A new list item is initialized with its state set to "TODO" (empty check box). Creating a new list item with another state can be done in the following way:

    li.Add("Sub item 1").state = gora.

Currently, only "TODO" (empty checkbox) and "DONE" (checkbox with "X") are supported. More will be added later on.
Other operations that are supported include:

    li.MoveUp(n) // Moves sub item n up one step
    li.MoveDown(n) // Moves sub item n down one step
    li.Delete(n) // Deletes sub item n

Finally, writing your list out to a readable text format is done in the following way:

    bytes, err := l.MarshalText()

At present, MarshalText() will not return any errors, but has return values ([]byte, error) to comply with the TextMarhsaler interface.

Example of TODO file
====================
- [ ] Buy groceries
  - [ ] Milk
  - [ ] Eggs
  - [X] Stuff

(In the future: ordered lists)
- [ ] Develop awesome function
   1. [X] Specify
   2. [X] Write test
   3. [ ] Implement
   4. [ ] Test

Notes for the future
====================
- Numbered lists
- Other attributes of lists
    * List name
    * Deadline
- Other attributes of tasks
    * Deadline
    * Priority
    * More states
- Kanban "view"
- Time tracking
    * Summary of time tracked
