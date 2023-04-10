terraform {
  required_providers {
    beget = {
      source  = "hashicorp.com/edu/beget"
    }
  }
}

provider "beget" {
  token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcklkIjoiMTIxNTQ3NCIsImN1c3RvbWVyTG9naW4iOiJidW5ueXBlYWNlIiwiZW52Ijoid2ViIiwiZXhwIjoxNzExNjE3ODYxLCJpYXQiOjE2ODAwNjM4MDEsImlwIjoiMTg1LjIwMy4yMzYuMjMwIiwiaXNzIjoiYXV0aC5iZWdldC5jb20iLCJqdGkiOiIzYjhlMWIzZDc0NTk4ZmQyYWU0MGJmMWM3YzY1NDgzMiIsInBhcmVudExvZ2luIjoiIiwic3ViIjoiY3VzdG9tZXIifQ.f-O8HGeMw0bkqXuZFzztUtSNhHSunpGdemu6gZY--3Q8RqglmYRCJNoF-1oVM7hvFynBV2iiFXLMsxdUtPC6aIL0tQxn8ovHsjEzbPmpAE_cCk1Tl7dIVm7Eq901L91KY522W0sDE3lqRE_USof1N_ssn-V--zAdBOEjgrUGjZla25KRFIaKB3u728nH0cW9INK3NrTh5lr6QQzF8JYqzn2bPrN0jWhWQjqtw4sDULo9_O7VEp272kjgqQfbpGl-IYdaaDtIlXhoapaU3XpE_Qd1pROyUK6BiCHstKKRca5i5aqh49Zm3zXWYNNweFDBUPf2FwZR3lBj3nV182Y5ZQ"
}

resource "beget_server" "example" {

}

output "example_server" {
  value = beget_server.example
}
