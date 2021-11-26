# MIHP

MIHP is an abbreviation of "MIHP Is Http Probe"

// TODO : work on the cel-go custom functions.


## MIHP Minion

### MIHP Minion Group

Multiple Minion instance is considered if they're belong to 1 network IP group.
For example, 
- Each minion will run on a single IP address.
- If a minion have ip of `1.2.3.4` and another minion have ip of `1.2.3.5` then they belong to the same group the `1.2.3.x` group.
- Minion that newly join the network e.g. `1.2.3.7` should know it belong to the `1.2.3.x` group.
- So there's possible of 256 minnion (from `x.x.x.1` to `x.x.x.255`) in a single group.
- Each group can only have 1 group leader (the one that have the highest `ElectScore`) which will be re elected if that group leader can not be reached.

### Minion Group Leadership, Election and ElectScore

When a new Minion is created and running, it will be assigned with random `uint64` value what
we call as `ElectScore`. This `ElectScore` determine a minion rank in the group. The minion with higher `ElectScore` should
be deemed as the Group Leader. All other minion instance should automaticaly appoint those with highest `ElectScore` as `GroupLeader`.
All other minion instance is a `GroupMember`.

GroupLeader have the following tasks :

1. Hourly Retrieve Probing requests from Central. All probe that assigned for its networks.
2. Distribute Probing requests received form Central to all `GroupMembers` including for it self.
3. Receive Hourly Probe result from all member 
4. Carry-on all Probing activity assigned for its own.
5. Respond to PING request from member.

GroupMember have the following task;

1. Hourly Retrieve Probing Request from GroupLeader. In case of the newly appointed GroupLeader, the member could ask the new leader ASAP.
2. Carry-on all Probing activity assigned for its own.
3. Report Probe Result to the GroupLeader.
4. Regularly PONG the GroupLeader for activity.

If a GroupLeader become unreachable (e.g. downed), a new GroupLeader is appointed from one of the GroupMember through election process.
The following are the events when Minion Group Leader need to be elected:

1. No leader are known between the members of Minion in the group. E.g. the network is freshly started. And the 1st minion is online.
2. The known leader can not be reached. Lost PING or can not TCP connect to it.
3. A new member entered the group with higher `ElectScore`


