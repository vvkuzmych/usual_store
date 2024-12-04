-- Remove the inserted widgets in the 'down' migration
DELETE FROM widgets WHERE name = 'Widget' AND description = 'A very nice widget.';
DELETE FROM widgets WHERE name = 'Golden Plan' AND description = 'Discount 30% for more than 3 subscriptions';
