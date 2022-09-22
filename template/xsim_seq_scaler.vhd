library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

package scaler_pack is
  constant scale_int : integer := {{.ScaleN}};
  constant word_size : integer := {{.BitSize}};
  constant output_size : integer := {{.OutputSize}};
  type inputarray is array (0 to scale_int-1) of std_logic_vector(word_size-1 downto 0);
  type outputarray is array (0 to scale_int-1) of std_logic_vector(output_size-1 downto 0);
  end package scaler_pack;


library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use std.textio.all;
use ieee.std_logic_textio.all; 
use work.scaler_pack.all;

entity sim is
generic (word_size: integer:={{.BitSize}}; output_size: integer:={{.OutputSize}}); 
end sim;

architecture Behavioral of sim is
    signal clk, rst : STD_LOGIC;
    signal a, b : inputarray;
    signal output : outputarray; 
    signal test : STD_LOGIC_VECTOR (output_size-1 downto 0); 
    -- buffer for storing the text from input read-file
    file input_buf : text;  -- text is keyword
begin

--This if-structure within the template is used to change the testbench from pre to post synthesis for post-synth functional analysis
--Vivado unrolls the array defined in the package into individual ports!
--Xilinx forums tell you to rewrite your testbench manually, this is a best-effort approach to keep it automatic (might need manual changes if port names change) 
{{if .PostSim}}
testmod : entity work.{{.TopEntityName}} port map(clk=>clk, rst=>rst,
{{range $i := N .ScaleN}}\a[{{$i}}]\=>a({{$i}}), {{end}}
{{range $i := N .ScaleN}}\b[{{$i}}]\=>b({{$i}}), {{end}}
{{$first := true}}{{range $i := N .ScaleN}}{{if $first}}{{$first = false}}{{else}},{{end}}\prod[{{$i}}]\=>output({{$i}}){{end}}
);   
{{else}}
testmod : entity work.{{.TopEntityName}} port map(clk=>clk, rst=>rst, a=>a, b=>b, prod=>output);   
{{end}}

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
          for I in 0 to scale_int-1 loop
            a(I) <= val_col1;
            b(I) <= val_col2;
          end loop;
          test <= val_col3;

          wait for 10 ns;  --  to display results for 20 ns

          clk <= '1';

          wait for 10 ns;
          
          for I in 0 to scale_int-1 loop
            assert(output(I) = test) report "!!ERROR!!PATTERN!!" severity ERROR;
          end loop;        
          
        end loop;

        file_close(input_buf);      
        
        report "Simulation Completed" severity NOTE;
        wait;
    end process;


end Behavioral;
