package server

var content = "Duis laoreet consequat fermentum. Sed finibus tempor sapien sit amet sollicitudin. Suspendisse vestibulum lacus imperdiet arcu venenatis, non convallis neque faucibus. Donec hendrerit ultricies enim vitae pulvinar. Phasellus ullamcorper ultrices dui, ut imperdiet ipsum pretium a. Curabitur laoreet, lectus nec laoreet convallis, erat tellus rhoncus enim, nec rutrum lorem enim quis ipsum. Sed eleifend commodo orci ac posuere. Etiam commodo massa nisl, ac bibendum felis imperdiet non. Fusce purus ligula, faucibus at felis ut, consequat congue urna. Sed rutrum urna in erat dignissim consectetur."

type Summary struct {
	Id      int
	Avatar  string
	Name    string
	Image   string
	Content string
}

var Characters = []Summary{
	{1, "http://lorempixel.com/50/50/?id=1", "one", "http://lorempixel.com/400/400/?id=1", content},
	{2, "http://lorempixel.com/50/50/?id=2", "two", "http://lorempixel.com/400/400/?id=2", content},
	{3, "http://lorempixel.com/50/50/?id=3", "three", "http://lorempixel.com/400/400/?id=3", content},
	{4, "http://lorempixel.com/50/50/?id=4", "four", "http://lorempixel.com/400/400/?id=4", content},
}

var Planets = []Summary{
	{1, "http://lorempixel.com/50/50/?id=1", "p-1", "http://lorempixel.com/400/400/?id=1", content},
	{2, "http://lorempixel.com/50/50/?id=2", "p-2", "http://lorempixel.com/400/400/?id=2", content},
	{3, "http://lorempixel.com/50/50/?id=3", "p-3", "http://lorempixel.com/400/400/?id=3", content},
	{4, "http://lorempixel.com/50/50/?id=4", "p-4", "http://lorempixel.com/400/400/?id=4", content},
}

var Starships = []Summary{
	{1, "http://lorempixel.com/50/50/?id=1", "s-1", "http://lorempixel.com/400/400/?id=1", content},
	{2, "http://lorempixel.com/50/50/?id=2", "s-2", "http://lorempixel.com/400/400/?id=2", content},
	{3, "http://lorempixel.com/50/50/?id=3", "s-3", "http://lorempixel.com/400/400/?id=3", content},
	{4, "http://lorempixel.com/50/50/?id=4", "s-4", "http://lorempixel.com/400/400/?id=4", content},
}
