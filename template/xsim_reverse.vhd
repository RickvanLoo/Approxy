library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use std.textio.all;
use ieee.std_logic_textio.all; 
use ieee.numeric_std.all;

entity sim is
generic (word_size: integer:={{.BitSize}}); 
end sim;

architecture Behavioral of sim is
    signal a, b : STD_LOGIC_VECTOR (word_size-1 downto 0);
    signal output : STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0); 
    signal test : STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0); 
    -- buffer for storing the text from input read-file
    file input_buf : text;  -- text is keyword
    file output_buf : text;
begin

testmod : entity work.{{.TopEntityName}} port map(a=>a, b=>b, prod=>output);   

testbench : process
    variable read_col_from_input_buf : line; -- read lines one by one from input_buf
    variable output_line : line;
    variable val_col1, val_col2 : STD_LOGIC_VECTOR (word_size-1 downto 0); -- to save col1 and col2 values of 1 bit
    variable val_SPACE : character;  -- for spaces between data in file
    variable output_data_buffer : STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0);
    begin
        file_open(input_buf, "{{.TestFile}}",  read_mode); 
        file_open(output_buf, "out_{{.TestFile}}", write_mode);
        while not endfile(input_buf) loop
          readline(input_buf, read_col_from_input_buf);
          read(read_col_from_input_buf, val_col1);
          read(read_col_from_input_buf, val_SPACE);           -- read in the space character
          read(read_col_from_input_buf, val_col2);


          -- Pass the read values to signals
          a <= val_col1;
          b <= val_col2;
          wait for 20 ns;  --  to display results for 20 ns

          output_data_buffer := output;
          write(output_line,to_integer(unsigned(output_data_buffer)));
          writeline(output_buf, output_line);
          
          --assert(output = test) report "!!ERROR!!PATTERN!!" severity ERROR;
        end loop;

        file_close(input_buf);      
        file_close(output_buf);
        report "Simulation Completed" severity NOTE;
        wait;
    end process;


end Behavioral;
