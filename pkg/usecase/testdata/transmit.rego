package transmit

import rego.v1

slack contains {
    "channel": "#general",
} if {
    input.schema == "for_slack"
}
