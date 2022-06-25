node_name = "consul-client"
data_dir = "/tmp/consul/data"
ports {
    https = 8501
}
auto_config {
    enabled = true,
    intro_token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjMwNDczMDJjLWI0OGMtMjRjMC1lOTljLTdkYjBlMTc4YjI2YyJ9.eyJhdWQiOiJjb25zdWwtY2x1c3Rlci1kYzEiLCJjb25zdWwiOnsiaG9zdG5hbWUiOiJjb25zdWwtY2xpZW50In0sImV4cCI6MTY1NDM2NTM4NSwiaWF0IjoxNjU0MzIyMTg1LCJpc3MiOiIvdjEvaWRlbnRpdHkvb2lkYyIsIm5hbWVzcGFjZSI6InJvb3QiLCJzdWIiOiJmMDQ5NzljYS03ZTQ1LWE4Y2UtMDNmOS0yOTA3NjU2MGI3MTgifQ.clBHQl4YFLeWvVUBcU_8_7EiRBA6WrYVBNV_Qx9EXhoJ7DMOaQwZ88wyjae3_yBIpHs_dwJvTeTQ0WubBEC5TVoHH-CyfwAEOrelyUym1EQF9xQnbOYvA7zV1sV46OgfppY_MmwGsUGtUVF1khVADyTnNvhR2iHo0taWoLBRPKo3NaxjKeb5T4kyMxfkkGzJ6pIJLJX1eOo6EK_nc2RBuboChP-GxHHXIxURMJ8GKltlCAnI4TdPfeTxyvcLNWdXsarx1-bDnX4GGE34GDNKBrWx8Aql0dKXQvuS-BVO6G1xg8KodISBB8kyTeZP4ue-87KC6gtR2xDQpMbRUi3OUQ"
    server_addresses = [
        "sd.r3t.io"
    ]
}
tls {
    defaults {
        verify_incoming = true
        verify_outgoing = true
        ca_file = "/tmp/consul/agent.ca.pem"
    }
}
