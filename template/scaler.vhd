library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

package scaler_pack is
constant scale_int : integer := {{.ScaleN}};
constant word_size : integer := {{.BitSize}};
type inputarray is array (0 to scale_int-1) of std_logic_vector(word_size-1 downto 0);
type outputarray is array (0 to scale_int-1) of std_logic_vector(word_size *2 - 1 downto 0);
end package scaler_pack;


library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use work.scaler_pack.all;

entity {{.EntityName}} is
Port ( 
a : in  inputarray;
b : in  inputarray;
prod: out outputarray);
end {{.EntityName}};

architecture Behavioral of {{.EntityName}} is
begin

scale_add: for I in 0 to scale_int-1 generate
mult: entity work.{{.LUTName}} port map(a=>a(I), b=>b(I),prod=>prod(I));
end generate scale_add;


end Behavioral;