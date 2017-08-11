provider "wavefront" {
  address = "spaceape.wavefront.com"
}

resource "wavefront_alert" "test_alert" {
  name = "Terraform Test Alert"
  target = "test@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=flamingo-int and cpu=cpu-total and service=game-service) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=flamingo-int and cpu=cpu-total and service=game-service)"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
    "flamingo"
  ]
}

resource "wavefront_alert" "test_alert1" {
  name = "Test Alert 1"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
    "test"
  ]
}
resource "wavefront_alert" "test_alert2" {
  name = "Test Alert 2"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
    "test"
  ]
}
resource "wavefront_alert" "test_alert3" {
  name = "Test Alert 3"
  target = "terraform@example.com"
  condition = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  display_expression = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes = 5
  resolve_after_minutes = 5
  severity = "WARN"
  tags = [
    "terraform",
    "test"
  ]
}