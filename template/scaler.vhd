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
use work.scaler_pack.all;

entity {{.EntityName}} is
Port (
{{if .MAC}}
clk : in std_logic;
rst : in std_logic;
{{end}}     
a : in  inputarray;
b : in  inputarray;
prod: out outputarray);
end {{.EntityName}};

architecture Behavioral of {{.EntityName}} is
begin

scale_add: for I in 0 to scale_int-1 generate
{{if .MAC}}
mult: entity work.{{.LUTName}} port map(clk=>clk, rst=>rst, a=>a(I), b=>b(I),prod=>prod(I));
{{else}}
mult: entity work.{{.LUTName}} port map(a=>a(I), b=>b(I),prod=>prod(I));
{{end}}
end generate scale_add;


end Behavioral;