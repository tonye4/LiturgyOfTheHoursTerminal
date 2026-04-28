# TODOS
- [ ] Make proper usage docs
- [ ] Add tabs for viewing certain prayers
- [x] Add text wrapping (lines split on elements)
    - Note: <s>I believe the solution is going through each element and splitting the line where there's periods.</s>
    - Note: The splitting was added via regex pattern matches on punctuation to add more line breaks, however
    there are still strings that are far too long and I can't seem to find a nice way to format these strings.
    Perhaps a more idiomatic solution would be to convert our blocks of text into markdown and redner them out that way.
    The ![Circumflex repo](https://github.com/bensadeh/circumflex) does a good job at parsing text from hackernews.
    Perhaps the way forward is implementing a parser similar to circumflex.
- [ ] Remove inline css from prayers
- [ ] Add command entry point to run the source
- [x] Document inline 
- [ ] Not all prayers start from 0% within the paginator. Fix this bug.
- [ ] On keywords such as PSALMODY, Ribbon placement...highlight red and center align
- [ ] Center align viewport
- [ ] Potentially implement a character limit to break long continuous blocks of text. 
When capturing text, instead of making a single large string with the string builder, perhaps make it so each time a line is read, store that information into a relational data structure (eg; dict) so that each line can be accessed individually.
This is useful because you can treat each line as an object which makes formatting easier; eg: for character counting.
    - Alternative: When pulling each text node from the tree, have each string's word count counted and if 
    it's above some arbitrary word count set, then append a `\n` in the middle of the string. 
    (Let's say 15 words in a sentence would indicate the need for a line break.)
