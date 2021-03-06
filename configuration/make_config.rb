# frozen_string_literal: true
require 'yaml'

DEFAULT_ARGS = {
    in: File.join(__dir__, 'templates'),
    out: File.join(__dir__, 'out'),
    serials: File.join(__dir__, 'serials.yml')
}.freeze
USAGE = <<~TEXT
    usage: make_config [--in=./templates] [--out=./out] [--serials=serials.yml]
TEXT

def symbolize_keys(hash)
    hash.transform_keys(&:to_sym)
end

# parse an array of strings of the format "--key=value" or "-flag" and return
# a hash of `:key => value` or `:flag => nil` pairs.
def parse_args(arr)
    symbolize_keys(arr.flat_map { |s| s.scan(/--?([^=\s]+)(?:=(\S+))?/) }.to_h)
end

# string -> hash, bool
def read_serials(filename)
    result = YAML.load(File.read(filename))
    ok = result.is_a?(Hash) && !result.empty?

    return result, ok
end

# string, hash -> string
def replace_placeholders(filename, serials)
    in_contents = File.read(filename)

    serials.each do |k, v|
        in_contents.gsub!(":#{k}:", v)
    end

    in_contents
end

# string, string -> int
def write_outfile(file, contents)
    file.write(contents)
end

def verify_plist(filename)
    system("plutil -s -- #{filename}")
end

def exit_with(code, message = nil, file = STDERR)
    file.puts(message) if message

    exit(code)
end

def main
    exit_with(0, USAGE, STDOUT) if ARGV[0] == '--help'

    args = DEFAULT_ARGS.merge(parse_args(ARGV))
    in_directory = args[:in]
    out_directory = args[:out]
    serials_filename = args[:serials]

    # Ensure the provided directories can be read from and written to:
    ok = File.directory?(in_directory) && File.readable?(in_directory)
    exit_with(2, "Input directory #{in_directory} is not readable") unless ok

    ok = File.directory?(out_directory) && File.writable?(out_directory)
    exit_with(2, "Output directory #{out_directory} is not writable") unless ok

    # Ensure the provided serials file exists, is readable, and contains
    # usable data:
    ok = File.file?(serials_filename) && File.readable?(serials_filename)
    exit_with(3, "File #{serials_filename} is not readable") unless ok

    serials, ok = read_serials(serials_filename)
    exit_with(4, "File #{serials_filename} contains no values") unless serials

    templates = Dir.glob(File.join(in_directory, '*.plist'))
    names = templates.map { |s| File.basename(s) }.join(', ')
    puts "Processing #{names}"
    all_ok = true

    templates.each do |template|
        in_file = File.open(template, File::RDONLY)
        out_filename = File.join(out_directory, File.basename(in_file))
        len = 0

        File.open(out_filename, 'w') do |file|
            out_contents = replace_placeholders(template, serials)
            len = write_outfile(file, out_contents)
        end

        ok = verify_plist(out_filename)

        puts "#{out_filename}: OK, #{len} bytes" if ok

        all_ok &&= ok
    end

    exit_with(all_ok ? 0 : 5)
end

main
