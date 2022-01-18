library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
use ieee.numeric_std.all;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}); 
Port ( 
--Generic
clk : in std_logic;
rst : in std_logic;
--Input Signal: Y=input, S=reference
ReY : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
ImY : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
ReS : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
ImS : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
-- STD DEV
dReYReS : in STD_LOGIC_VECTOR (word_size-1 downto 0);
dImYImS : in STD_LOGIC_VECTOR (word_size-1 downto 0);
dReYImS : in STD_LOGIC_VECTOR (word_size-1 downto 0);
dImyReS : in STD_LOGIC_VECTOR (word_size-1 downto 0);
--Output
Real: out STD_LOGIC_VECTOR (word_size * 3 - 1 downto 0);
Img: out STD_LOGIC_VECTOR (word_size * 3 - 1 downto 0));
end {{.EntityName}};
    
architecture Behavioral of {{.EntityName}} is
    signal ReYReS, ImYImS, ReYImS, ImYReS :  STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0) := (others => '0');
begin

ReYReS_mac: entity work.{{ (index .MACArray 0).EntityName }} port map(clk=>clk, rst=>rst, a=>ReY, b=>ReS,prod=>ReYReS);
ImYImS_mac: entity work.{{ (index .MACArray 1).EntityName }} port map(clk=>clk, rst=>rst, a=>ImY, b=>ImS,prod=>ImYImS);
ReYImS_mac: entity work.{{ (index .MACArray 2).EntityName }} port map(clk=>clk, rst=>rst, a=>ReY, b=>ImS,prod=>ReYImS);
ImYReS_mac: entity work.{{ (index .MACArray 3).EntityName }} port map(clk=>clk, rst=>rst, a=>ImY, b=>ReS,prod=>ImYReS);


process(ReYReS, ImYImS, ReYImS, ImYReS, dReYReS, dImYImS, dReYImS, dImYReS)
begin
    Real <= std_logic_vector((unsigned(ReYReS)*unsigned(dReYReS))+(unsigned(ImYImS)*unsigned(dImYImS)));
    Img <= std_logic_vector((unsigned(ReYImS)*unsigned(dReYImS))+(unsigned(ImYReS)*unsigned(dImYReS)));
end process;


end Behavioral;