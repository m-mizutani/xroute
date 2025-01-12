package transmit

import rego.v1

slack contains {
    "title": "Hello",
    "emoji": ":wave:",
    "channel": "#github-notify",
    "body": input.data,
} if {
    is_string(input.data)
}

slack contains {
    "title": "Hello",
    "emoji": ":wave:",
    "channel": "#github-notify",
    "body": json.marshal(input.data),
} if {
    is_object(input.data)
}
