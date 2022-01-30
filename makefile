seeder:
	if [ ! -d $(path) ]; then mkdir $(path); fi
	cp ./seed/template.json $(path)/$(name).json
