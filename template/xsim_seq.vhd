library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use std.textio.all;
use ieee.std_logic_textio.all; 

entity sim is
generic (word_size: integer:={{.BitSize}}; output_size: integer:={{.OutputSize}}); 
end sim;

architecture Behavioral of sim is
    signal clk, rst : STD_LOGIC;
    signal a, b : STD_LOGIC_VECTOR (word_size-1 downto 0);
    signal output : STD_LOGIC_VECTOR (output_size-1 downto 0); 
    signal test : STD_LOGIC_VECTOR (output_size-1 downto 0); 
    -- buffer for storing the text from input read-file
    file input_buf : text;  -- text is keyword
begin

testmod : entity work.{{.TopEntityName}} port map(clk=>clk, rst=>rst, a=>a, b=>b, prod=>output);   

testbench : process
    variable read_col_from_input_buf : line; -- read lines one by one from input_buf
    variable val_col1, val_col2 : STD_LOGIC_VECTOR (word_size-1 downto 0); -- to save col1 and col2 values of 1 bit
    variable val_col3 : STD_LOGIC_VECTOR (output_size-1 downto 0); -- to save col3 value of 2 bit
    variable val_SPACE : character;  -- for spaces between data in file
    begin
        rst <= '1';

        wait for 20 ns;

        rst <= '0';

        wait for 20 ns;


        file_open(input_buf, "{{.TestFile}}",  read_mode); 
        while not endfile(input_buf) loop
          readline(input_buf, read_col_from_input_buf);
          read(read_col_from_input_buf, val_col1);
          read(read_col_from_input_buf, val_SPACE);           -- read in the space character
          read(read_col_from_input_buf, val_col2);
          read(read_col_from_input_buf, val_SPACE);           -- read in the space character
          read(read_col_from_input_buf, val_col3);

          -- Pass the read values to signals
          clk <= '0';
          a <= val_col1;
          b <= val_col2;
          test <= val_col3;

          wait for 20 ns;  --  to display results for 20 ns

          clk <= '1';

          wait for 20 ns;
          
          assert(output = test) report "!!ERROR!!PATTERN!!" severity ERROR;
        end loop;

        file_close(input_buf);      
        
        report "Simulation Completed" severity NOTE;
        wait;
    end process;


end Behavioral;
