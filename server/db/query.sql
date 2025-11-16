-- name: LinkWithIdExists :one
SELECT EXISTS (
    SELECT 1
    FROM links
    WHERE id = $1
);

-- name: CreateLink :exec
INSERT INTO links (
    id, url
) VALUES (
    $1, $2
);

-- name: GetLinkById :one
SELECT * FROM links
WHERE id = $1
LIMIT 1;

-- name: CreateClick :exec
INSERT INTO clicks (
    link_id, referrer
) VALUES (
    $1, $2
);

-- name: CountRecentClicks :one
SELECT COUNT(*) FROM clicks
WHERE link_id = $1 AND created_at >= NOW() - INTERVAL '3 days';

-- name: GetWeeklyClickStats :many
SELECT referrer, COUNT(*) FROM clicks
WHERE link_id = $1 AND created_at >= NOW() - INTERVAL '7 days'
GROUP BY referrer;