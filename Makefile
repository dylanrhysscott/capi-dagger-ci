infra:
	dagger call deploy-infra \
		--path "./infra" \
		--token=env:DIGITALOCEAN_ACCESS_TOKEN \
		--spaces-access-key=env:AWS_ACCESS_KEY_ID \
		--spaces-secret-key=env:AWS_SECRET_ACCESS_KEY 

capi:
	dagger call deploy-infra \
		--token=env:DIGITALOCEAN_ACCESS_TOKEN \
		--cluster-name="dscott-capi"