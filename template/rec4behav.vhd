library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use ieee.numeric_std.all;

entity {{.EntityName}} is
generic (word_size: integer:=4); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};
    
architecture Behavioral of {{.EntityName}} is
    signal AH_BH_OUT, AH_BL_OUT, AL_BH_OUT, AL_BL_OUT :  std_logic_vector(7 downto 0) := (others => '0');
    signal AH, AL, BH, BL : std_logic_vector(1 downto 0) := (others => '0');
    signal factor1, factor2, factor3, factor4 : unsigned(7 downto 0) := (others => '0');
begin

AH_BH: entity work.AH_BH port map(a=>AH, b=>BH,prod=>AH_BH_OUT(3 downto 0));
AH_BL: entity work.AH_BL port map(a=>AH, b=>BL,prod=>AH_BL_OUT(3 downto 0));
AL_BH: entity work.AL_BH port map(a=>AL, b=>BH,prod=>AL_BH_OUT(3 downto 0));
AL_BL: entity work.AL_BL port map(a=>AL, b=>BL,prod=>AL_BL_OUT(3 downto 0));

AH <= A(3 downto 2);
AL <= A(1 downto 0);
BH <= B(3 downto 2);
BL <= B(1 downto 0);

factor1 <= unsigned(AL_BL_OUT);
factor2 <= shift_left(unsigned(AL_BH_OUT),2);
factor3 <= shift_left(unsigned(AH_BL_OUT),2);
factor4 <= shift_left(unsigned(AH_BH_OUT),4);

prod <= std_logic_vector(factor1 + factor2 + factor3 + factor4);

end Behavioral;