function load_lines_from_file(file)
    lines = {}

    local f=io.open(file,"r")
    if f~=nil then
        io.close(f)
    else
        return lines
    end

    for line in io.lines(file) do
        if not (line == '') then
            lines[#lines + 1] = line
        end
    end

    return lines
end

urls = load_lines_from_file("random.urls")

wrk.headers['X-Auth-Token'] = '[AUTH_TOKEN]'

counter = 0
request = function()
    -- Get the next array element
    url_path = urls[counter]

    counter = counter + 1

    -- If the counter is longer than the urls array length then reset it
    if counter > #urls then
        counter = 0
    end

    -- Return the request object with the current URL path
    return wrk.format('GET', url_path, nil, nil)
end
