# Meeting Cost Clock
A clock that tells you how much a meeting is costing you, written in Go with Moria.

## Do you have meeting-itis?

Are you tired of having too many meetings, with too many people? Are you fed up with meetings constantly running over time?

Then this app might just be the tool you need to passive-agressively remind everyone that money is time, and that time 
flies quickly.

You can begin a meeting by entering a duration, adding some attendees, and then hitting *Start the meeting*:
<figure>
<figcaption><label for="begin-image"><b>Fig 1: Beginning a meeting</b></label></figcaption><br>
<img id="begin-image" src="documentation-assets/beginning-a-meeting.png" width=600></img>
</figure>

This takes you to the following dramatic scene where everyone can see the time running out!
<figure>
<figcaption><label for="during-image"><b>Fig 2: A meeting in progress</b></label></figcaption><br>
<img src="documentation-assets/meeting-in-progress.png" width=600></img>
</figure>

What's more, once the meeting has run over time, the cost will turn an angry red, to let everyone 
know that precious time is being wasted.

> But doesn't this already exist?

Yeah, you got me. After I got the idea to make this, I found out that a number of similar tools 
built around estimating meeting costs already exist. So why did I decide to build another one myself?

1. The first reason is it's such a simple idea to reimplement that it hardly makes a difference that
something like this already exists. Just like Classical painters all drawing the Last Supper and leaving 
their individual mark on it. 

2. The second, and larger reason, is that I wanted a chance to use and showcase the Moria framework that 
I built as part of [go-mithril](https://github.com/danverbraganza/go-mithril). Using Moria, this 
Single-Page Web App was written entirely in Go, and compiled to Javascript using GopherJS.


