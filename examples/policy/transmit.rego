package route

import rego.v1

slack contains {
    "title": "Hello",
    "emoji": ":wave:",
    "channel": "#github-notify",
    "body": input.body,
} if {
    is_string(input.body)
}

slack contains {
    "title": "Hello",
    "emoji": ":wave:",
    "channel": "#github-notify",
    "body": json.marshal(input.body),
} if {
    is_object(input.body)
}
